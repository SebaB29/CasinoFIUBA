package ruleta

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	repositoriesJuegos "casino/repositories/juegos"
	"time"
)

const (
	TipoTransaccionApuesta  = "apuesta"
	TipoTransaccionGanancia = "ganancia"
)

type RuletaService struct {
	usuarioRepository     repositories.UsuarioRepositoryInterface
	jugadaRepository      repositoriesJuegos.JugadaRuletaRepositoryInterface
	transaccionRepository repositories.TransaccionRepositoryInterface
}

func NewRuletaService() *RuletaService {
	return &RuletaService{
		usuarioRepository:     repositories.NewUsuarioRepository(db.DB),
		jugadaRepository:      repositoriesJuegos.NewJugadaRuletaRepository(db.DB),
		transaccionRepository: repositories.NewTransaccionRepository(db.DB),
	}
}

func (ruletaService *RuletaService) Jugar(usuarioID uint, jugada dto.RuletaRequestDTO) (dto.RuletaResponseDTO, error) {
	usuario, err := ruletaService.ValidarApuesta(usuarioID, jugada)
	if err != nil {
		return dto.RuletaResponseDTO{}, err
	}

	resultado := EjecutarRuleta(jugada)

	if err := ruletaService.procesarResultado(usuario, resultado); err != nil {
		return dto.RuletaResponseDTO{}, err
	}

	if err := ruletaService.registrarJugada(usuario.ID, jugada, resultado); err != nil {
		return dto.RuletaResponseDTO{}, err
	}

	return resultado, nil
}

func (ruletaService *RuletaService) ValidarApuesta(usuarioID uint, jugada dto.RuletaRequestDTO) (*models.Usuario, error) {
	const MontoMinimo = 1.0

	usuario, err := ruletaService.usuarioRepository.ObtenerPorID(usuarioID)
	if err != nil || usuario == nil {
		return nil, errores.ErrUsuarioNoEncontrado
	}

	if err := ValidarJugada(jugada); err != nil {
		return nil, err
	}

	if jugada.Monto < MontoMinimo {
		return nil, errores.ErrMontoInsuficiente
	}

	if usuario.Saldo < jugada.Monto {
		return nil, errores.ErrSaldoInsuficiente
	}

	return usuario, nil
}

func (ruletaService *RuletaService) procesarResultado(usuario *models.Usuario, resultado dto.RuletaResponseDTO) error {
	// Registrar transacción de apuesta
	if err := ruletaService.registrarTransaccion(usuario.ID, TipoTransaccionApuesta, resultado.MontoApostado); err != nil {
		return err
	}

	// Actualizar saldo
	usuario.Saldo = usuario.Saldo - resultado.MontoApostado + resultado.Ganancia
	if err := ruletaService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	// Registrar transacción de ganancia
	if resultado.Ganancia > 0 {
		if err := ruletaService.registrarTransaccion(usuario.ID, TipoTransaccionGanancia, resultado.Ganancia); err != nil {
			return err
		}
	}

	return nil
}

func (ruletaService *RuletaService) registrarJugada(usuarioID uint, jugadaInicial dto.RuletaRequestDTO, resultado dto.RuletaResponseDTO) error {
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
		NumeroGanador: resultado.NumeroGanador,
		ColorGanador:  resultado.ColorGanador,
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
