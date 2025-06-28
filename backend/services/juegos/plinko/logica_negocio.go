package plinko

import (
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/models"
	"casino/utils"
	"encoding/json"
	"time"
)

func (plinkoService *PlinkoService) validarJugada(usuarioID uint, monto float64) (*models.Usuario, error) {
	const MontoMinimo = 1.0

	if monto < MontoMinimo {
		return nil, errores.ErrMontoInsuficiente
	}

	usuario, err := plinkoService.usuarioRepository.ObtenerPorID(usuarioID)
	if err != nil || usuario == nil {
		return nil, errores.ErrUsuarioNoEncontrado
	}

	if usuario.Saldo < monto {
		return nil, errores.ErrSaldoInsuficiente
	}

	return usuario, nil
}

func (plinkoService *PlinkoService) procesarResultado(usuario *models.Usuario, monto float64, resultado dto.PlinkoResponseDTO) error {
	// Registrar transacción de apuesta
	if err := plinkoService.registrarTransaccion(usuario.ID, utils.TipoTransaccionApuesta, monto); err != nil {
		return err
	}

	// Actualizar saldo
	usuario.Saldo = usuario.Saldo - monto + resultado.Ganancia
	if err := plinkoService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	// Registrar transacción de ganancia
	if err := plinkoService.registrarTransaccion(usuario.ID, utils.TipoTransaccionGanancia, resultado.Ganancia); err != nil {
		return err
	}

	return nil
}

func (plinkoService *PlinkoService) registrarTransaccion(usuarioID uint, tipoTransaccion string, monto float64) error {
	transaccion := &models.Transaccion{
		UsuarioID: usuarioID,
		Tipo:      tipoTransaccion,
		Monto:     monto,
	}

	return plinkoService.transaccionRepository.Crear(transaccion)
}

func (plinkoService *PlinkoService) registrarJugada(usuarioID uint, monto float64, resultado dto.PlinkoResponseDTO) error {
	trayectoJSON, err := json.Marshal(resultado.Trayecto)
	if err != nil {
		return err
	}

	jugada := &models.JugadaPlinko{
		UsuarioID:     usuarioID,
		MontoApostado: monto,
		Multiplicador: resultado.Multiplicador,
		Ganancia:      resultado.Ganancia,
		PosicionFinal: resultado.PosicionFinal,
		Trayecto:      string(trayectoJSON),
		Fecha:         time.Now(),
	}
	return plinkoService.jugadaRepository.Crear(jugada)
}
