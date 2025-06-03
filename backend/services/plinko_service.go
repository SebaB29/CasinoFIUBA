package services

import (
	"casino/db"
	"casino/dto"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	"math/rand"
	"time"
)

type PlinkoService struct {
	usuarioRepository repositories.UsuarioRepositoryInterface
	jugadaRepository  repositories.JugadaPlinkoRepositoryInterface
}

func NewPlinkoService() *PlinkoService {
	return &PlinkoService{
		usuarioRepository: repositories.NewUsuarioRepository(db.DB),
		jugadaRepository:  repositories.NewJugadaPlinkoRepository(db.DB),
	}
}

// Lógica del juego
func ejecutarPlinko(monto float64) dto.PlinkoResponseDTO {
	const niveles = 8

	// Crear generador local con semilla única
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Simular rebotes
	pos := 0
	for i := 0; i < niveles; i++ {
		direccion := r.Intn(2) // 0 = izquierda, 1 = derecha
		pos += direccion
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

	usuario.Saldo = usuario.Saldo - monto + resultado.Ganancia
	if err := plinkoService.usuarioRepository.Actualizar(usuario); err != nil {
		return dto.PlinkoResponseDTO{}, err
	}

	if err := plinkoService.guardarJugada(usuario.ID, monto, resultado); err != nil {
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

func (plinkoService *PlinkoService) guardarJugada(usuarioID uint, monto float64, resultado dto.PlinkoResponseDTO) error {
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
