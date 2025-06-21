package slot

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/repositories"
	repositoriesJuegos "casino/repositories/juegos"
)

type SlotService struct {
	usuarioRepository     repositories.UsuarioRepositoryInterface
	transaccionRepository repositories.TransaccionRepositoryInterface
	jugadaRepository      repositoriesJuegos.JugadaSlotRepositoryInterface
}

func NewSlotService() *SlotService {
	return &SlotService{
		usuarioRepository:     repositories.NewUsuarioRepository(db.DB),
		transaccionRepository: repositories.NewTransaccionRepository(db.DB),
		jugadaRepository:      repositoriesJuegos.NewJugadaSlotRepository(db.DB),
	}
}

func (slotService *SlotService) Jugar(userID uint, req dto.SlotRequestDTO) (dto.SlotResponseDTO, error) {
	usuario, err := slotService.validarJugada(userID, req.Monto)
	if err != nil {
		return dto.SlotResponseDTO{}, err
	}

	rondas, ganancia := generarResultadoSlot(req.Monto)

	if err := slotService.procesarResultado(usuario, req.Monto, ganancia); err != nil {
		return dto.SlotResponseDTO{}, err
	}

	if err := slotService.registrarJugada(usuario.ID, req.Monto, ganancia, rondas); err != nil {
		return dto.SlotResponseDTO{}, err
	}

	return dto.SlotResponseDTO{
		Rondas:   rondas,
		Ganancia: ganancia,
		Mensaje:  generarMensajeSlot(ganancia),
	}, nil
}
