package plinko

import (
	"casino/db"
	dto "casino/dto/juegos"
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
