package ruleta

import (
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/models"
	"casino/utils"
	"time"
)

func (ruletaService *RuletaService) ValidarApuesta(usuario *models.Usuario, jugada dto.RuletaRequestDTO) error {
	const MontoMinimo = 1.0

	if jugada.Monto < MontoMinimo {
		return errores.ErrMontoInsuficiente
	}

	if usuario.Saldo < jugada.Monto {
		return errores.ErrSaldoInsuficiente
	}

	// La función se encuentra implementada en validar_jugada.go
	if err := ValidarJugada(jugada); err != nil {
		return err
	}

	return nil
}

func (ruletaService *RuletaService) procesarResultado(usuario *models.Usuario, resultado ResultadoRuleta) error {
	// Registrar transacción de apuesta
	if err := ruletaService.registrarTransaccion(usuario.ID, utils.TipoTransaccionApuesta, resultado.MontoApostado); err != nil {
		return err
	}

	// Actualizar saldo
	usuario.Saldo = usuario.Saldo - resultado.MontoApostado + resultado.Ganancia
	if err := ruletaService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	// Registrar transacción de ganancia
	if resultado.Ganancia > 0 {
		if err := ruletaService.registrarTransaccion(usuario.ID, utils.TipoTransaccionGanancia, resultado.Ganancia); err != nil {
			return err
		}
	}

	return nil
}

func (ruletaService *RuletaService) registrarJugada(usuarioID uint, jugadaInicial dto.RuletaRequestDTO, resultado ResultadoRuleta) error {
	jugada := &models.JugadaRuleta{
		UsuarioID:     usuarioID,
		MontoApostado: jugadaInicial.Monto,
		Ganancia:      resultado.Ganancia,
		TipoApuesta:   jugadaInicial.TipoApuesta,
		Numeros:       models.IntSlice(jugadaInicial.Numeros),
		Docena:        jugadaInicial.Docena,
		Color:         jugadaInicial.Color,
		Paridad:       jugadaInicial.Paridad,
		AltoBajo:      jugadaInicial.AltoBajo,
		NumeroGanador: resultado.NumeroGanador.Valor,
		ColorGanador:  resultado.NumeroGanador.Color,
		Fecha:         time.Now(),
	}

	return ruletaService.jugadaRepository.Crear(jugada)
}

func (ruletaService *RuletaService) registrarTransaccion(usuarioID uint, tipoTransaccion string, monto float64) error {
	transaccion := &models.Transaccion{
		UsuarioID: usuarioID,
		Tipo:      tipoTransaccion,
		Monto:     monto,
	}

	return ruletaService.transaccionRepository.Crear(transaccion)
}