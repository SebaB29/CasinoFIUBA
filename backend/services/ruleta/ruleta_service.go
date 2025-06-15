package ruleta

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/repositories"
	repositoriesJuegos "casino/repositories/juegos"
	"time"
)

type RuletaService struct {
	ruletaManager         *RuletaManager
	usuarioRepository     repositories.UsuarioRepositoryInterface
	jugadaRepository      repositoriesJuegos.JugadaRuletaRepositoryInterface
	transaccionRepository repositories.TransaccionRepositoryInterface
}

func NewRuletaService() *RuletaService {
	return &RuletaService{
		ruletaManager:         &RuletaManager{},
		usuarioRepository:     repositories.NewUsuarioRepository(db.DB),
		jugadaRepository:      repositoriesJuegos.NewJugadaRuletaRepository(db.DB),
		transaccionRepository: repositories.NewTransaccionRepository(db.DB),
	}
}

func (ruletaService *RuletaService) Jugar(usuarioID uint, jugada dto.RuletaRequestDTO) (chan ResultadoRuleta, error) {
	usuario, err := ruletaService.usuarioRepository.ObtenerPorID(usuarioID)
	if err != nil || usuario == nil {
		return nil, errores.ErrUsuarioNoEncontrado
	}

	// La función se encuentra implementada en logica_negocio.go
	if err := ruletaService.ValidarApuesta(usuario, jugada); err != nil {
		return nil, err
	}

	resultadoChannel := make(chan ResultadoRuleta, 1) // Canal con buffer de 1 para evitar bloqueo
	ruletaService.iniciarTemporizador(usuarioID, jugada, resultadoChannel)

	return resultadoChannel, nil
}

func (ruletaService *RuletaService) EjecutarRuleta() {
	ruletaActual := &ruletaService.ruletaManager.JuegoActual

	ruletaActual.Mutex.Lock()
	jugadas := ruletaActual.Jugadas

	// Reinicio las Jugadas y Timer para la proxima ronda
	ruletaActual.Jugadas = nil
	ruletaActual.TimerActivo = false

	ruletaActual.Mutex.Unlock()

	// La función se encuentra implementada en logica_juego.go
	numeroGanador := obtenerNumeroGanador()

	for _, jugada := range jugadas {
		usuarioID := jugada.UsuarioID
		usuario, err := ruletaService.usuarioRepository.ObtenerPorID(usuarioID)
		if err != nil || usuario == nil {
			continue
		}

		jugadaDeUsuario := jugada.Apuesta

		// La función se encuentra implementada en logica_juego.go
		multiplicador := calcularMultiplicador(jugadaDeUsuario, numeroGanador)

		resultado := ResultadoRuleta{
			MontoApostado: jugadaDeUsuario.Monto,
			Ganancia:      jugadaDeUsuario.Monto * multiplicador,
			NumeroGanador: numeroGanador,
		}

		if err := ruletaService.procesarResultado(usuario, resultado); err != nil {
			continue
		}

		if err := ruletaService.registrarJugada(usuarioID, jugadaDeUsuario, resultado); err != nil {
			continue
		}

		jugada.Resultado <- resultado
		close(jugada.Resultado)
	}
}

func (ruletaService *RuletaService) iniciarTemporizador(usuarioID uint, jugada dto.RuletaRequestDTO, resultadoChannel chan ResultadoRuleta) {
	ruletaActual := &ruletaService.ruletaManager.JuegoActual

	ruletaActual.Mutex.Lock()

	timerActivo := ruletaActual.TimerActivo
	if !timerActivo {
		ruletaActual.TimerActivo = true
	}

	ruletaActual.Jugadas = append(ruletaActual.Jugadas, JugadaConUsuario{
		Apuesta:   jugada,
		UsuarioID: usuarioID,
		Resultado: resultadoChannel,
	})
	ruletaActual.Mutex.Unlock()

	if !timerActivo {
		go func() {
			time.Sleep(15 * time.Second)
			ruletaService.EjecutarRuleta()
		}()
	}
}
