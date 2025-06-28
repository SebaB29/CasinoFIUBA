package slot

import (
	"casino/errores"
	"casino/models"
	"casino/utils"
	"encoding/json"
	"time"
)

func (slotService *SlotService) validarJugada(userID uint, monto float64) (*models.Usuario, error) {
	const MontoMinimo = 1.0
	if monto < MontoMinimo {
		return nil, errores.ErrMontoInsuficiente
	}

	usuario, err := slotService.usuarioRepository.ObtenerPorID(userID)
	if err != nil || usuario == nil {
		return nil, errores.ErrUsuarioNoEncontrado
	}

	if usuario.Saldo < monto {
		return nil, errores.ErrSaldoInsuficiente
	}

	return usuario, nil
}

func (slotService *SlotService) procesarResultado(usuario *models.Usuario, monto, ganancia float64) error {
	if err := slotService.registrarTransaccion(usuario.ID, utils.TipoTransaccionApuesta, monto); err != nil {
		return err
	}

	usuario.Saldo = usuario.Saldo - monto + ganancia
	if err := slotService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	if ganancia > 0 {
		if err := slotService.registrarTransaccion(usuario.ID, utils.TipoTransaccionGanancia, ganancia); err != nil {
			return err
		}
	}

	return nil
}

func (slotService *SlotService) registrarTransaccion(usuarioID uint, tipo string, monto float64) error {
	transaccion := &models.Transaccion{
		UsuarioID: usuarioID,
		Tipo:      tipo,
		Monto:     monto,
	}
	return slotService.transaccionRepository.Crear(transaccion)
}

func (slotService *SlotService) registrarJugada(usuarioID uint, monto float64, ganancia float64, rondas [][]string) error {
	rondasJSON, err := json.Marshal(rondas)
	if err != nil {
		return err
	}

	jugada := &models.JugadaSlot{
		UsuarioID:     usuarioID,
		MontoApostado: monto,
		Ganancia:      ganancia,
		Rondas:        string(rondasJSON),
		Fecha:         time.Now(),
	}
	return slotService.jugadaRepository.Crear(jugada)
}
