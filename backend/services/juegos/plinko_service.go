package services

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	repositoriesJuegos "casino/repositories/juegos"
	"math/rand"
	"time"
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

// Lógica del juego
func ejecutarPlinko(monto float64) dto.PlinkoResponseDTO {
	const niveles = 8
	const centro = 4

	// Crear generador local con semilla única
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Posición inicial
	pos := centro

	// Simular rebotes
	for i := 0; i < niveles; i++ {
		direccion := r.Intn(2) // 0 = izquierda, 1 = derecha
		if direccion == 0 && pos > 0 {
			pos--
		} else if direccion == 1 && pos < 8 {
			pos++
		}
	}

	// Multiplicadores por posición
	multiplicadores := []float64{0.5, 0.8, 1.0, 1.4, 3.0, 1.4, 1.0, 0.8, 0.5}

	multiplicador := multiplicadores[pos]
	ganancia := monto * multiplicador

	return dto.PlinkoResponseDTO{
		PosicionFinal: pos,
		Multiplicador: multiplicador,
		Ganancia:      ganancia,
	}
}

func (plinkoService *PlinkoService) Jugar(usuarioID uint, monto float64) (dto.PlinkoResponseDTO, error) {
	usuario, err := plinkoService.validarJugada(usuarioID, monto)
	if err != nil {
		return dto.PlinkoResponseDTO{}, err
	}

	resultado := ejecutarPlinko(monto)

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
	if err := plinkoService.registrarTransaccion(usuario.ID, "apuesta", monto); err != nil {
		return err
	}

	// Actualizar saldo
	usuario.Saldo = usuario.Saldo - monto + resultado.Ganancia
	if err := plinkoService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	// Registrar transacción de ganancia
	if err := plinkoService.registrarTransaccion(usuario.ID, "ganancia", resultado.Ganancia); err != nil {
		return err
	}

	return nil
}

func (plinkoService *PlinkoService) registrarJugada(usuarioID uint, monto float64, resultado dto.PlinkoResponseDTO) error {
	jugada := &models.JugadaPlinko{
		UsuarioID:     usuarioID,
		MontoApostado: monto,
		Multiplicador: resultado.Multiplicador,
		Ganancia:      resultado.Ganancia,
		PosicionFinal: resultado.PosicionFinal,
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
