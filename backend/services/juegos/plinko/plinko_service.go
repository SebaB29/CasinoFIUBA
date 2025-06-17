package plinko

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	repositoriesJuegos "casino/repositories/juegos"
	"casino/utils"
	"encoding/json"
	"math/rand"
	"time"
)

const (
	TipoTransaccionApuesta  = "apuesta"
	TipoTransaccionGanancia = "ganancia"
)

type PlinkoService struct {
	usuarioRepository     repositories.UsuarioRepositoryInterface
	jugadaRepository      repositoriesJuegos.JugadaPlinkoRepositoryInterface
	transaccionRepository repositories.TransaccionRepositoryInterface
}

func NewPlinkoService() *PlinkoService {
	return &PlinkoService{
		usuarioRepository:     repositories.NewUsuarioRepository(db.DB),
		jugadaRepository:      repositoriesJuegos.NewJugadaPlinkoRepository(db.DB),
		transaccionRepository: repositories.NewTransaccionRepository(db.DB),
	}
}

func (plinkoService *PlinkoService) Jugar(usuarioID uint, monto float64) (dto.PlinkoResponseDTO, error) {
	usuario, err := plinkoService.validarJugada(usuarioID, monto)
	if err != nil {
		return dto.PlinkoResponseDTO{}, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	resultado := EjecutarPlinko(monto, r)

	if err := plinkoService.procesarResultado(usuario, monto, resultado); err != nil {
		return dto.PlinkoResponseDTO{}, err
	}

	if err := plinkoService.registrarJugada(usuario.ID, monto, resultado); err != nil {
		return dto.PlinkoResponseDTO{}, err
	}

	return resultado, nil
}

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

func (plinkoService *PlinkoService) registrarTransaccion(usuarioID uint, tipoTransaccion string, monto float64) error {
	transaccion := &models.Transaccion{
		UsuarioID: usuarioID,
		Tipo:      tipoTransaccion,
		Monto:     monto,
	}

	return plinkoService.transaccionRepository.Crear(transaccion)
}
