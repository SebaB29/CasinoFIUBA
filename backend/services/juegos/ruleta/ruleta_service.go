package ruleta

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/errores"
	"casino/repositories"
	repositoriesJuegos "casino/repositories/juegos"
	"time"

	"github.com/gorilla/websocket"
)

type RuletaService struct {
	ruleta                *RuletaEnJuego
	usuarioRepository     repositories.UsuarioRepositoryInterface
	jugadaRepository      repositoriesJuegos.JugadaRuletaRepositoryInterface
	transaccionRepository repositories.TransaccionRepositoryInterface
}

func NewRuletaService() *RuletaService {
	return &RuletaService{
		ruleta: &RuletaEnJuego{
			ConexionesWS: make(map[uint]*websocket.Conn),
			Estado:       EstadoEsperandoApuestas,
		},
		usuarioRepository:     repositories.NewUsuarioRepository(db.DB),
		jugadaRepository:      repositoriesJuegos.NewJugadaRuletaRepository(db.DB),
		transaccionRepository: repositories.NewTransaccionRepository(db.DB),
	}
}

func (ruletaService *RuletaService) Jugar(usuarioID uint, jugada dto.RuletaRequestDTO, conn *websocket.Conn) error {
	usuario, err := ruletaService.usuarioRepository.ObtenerPorID(usuarioID)
	if err != nil || usuario == nil {
		return errores.ErrUsuarioNoEncontrado
	}

	// La función se encuentra implementada en logica_negocio.go
	if err := ruletaService.ValidarApuesta(usuario, jugada); err != nil {
		return err
	}

	ruletaActual := ruletaService.ruleta

	ruletaActual.Mutex.Lock()

	if ruletaActual.Estado != EstadoEsperandoApuestas {
		ruletaActual.Mutex.Unlock()
		return errores.ErrApuestaFueraDeTiempo
	}

	// Guardamos la jugada y conexión
	ruletaActual.ConexionesWS[usuarioID] = conn
	ruletaActual.Jugadas = append(ruletaActual.Jugadas, JugadaConUsuario{
		Apuesta:   jugada,
		UsuarioID: usuarioID,
	})

	// Si no hay timer activo, lo arrancamos
	if !ruletaActual.TimerActivo {
		ruletaActual.TimerActivo = true
		ruletaActual.Estado = EstadoGirando
		go func() {
			time.Sleep(15 * time.Second)
			ruletaService.anunciarCierreApuestas(ruletaActual)
			time.Sleep(5 * time.Second)
			ruletaService.EjecutarRuleta()
		}()
	}

	ruletaActual.Mutex.Unlock()
	return nil
}

func (ruletaService *RuletaService) EjecutarRuleta() {
	ruletaActual := ruletaService.ruleta

	numeroGanador := obtenerNumeroGanador()
	resultados := ruletaService.evaluarJugadas(ruletaActual.Jugadas, numeroGanador)
	ruletaService.enviarResultados(ruletaActual, resultados, numeroGanador)

	ruletaService.prepararNuevaRonda(ruletaActual)
}

func (ruletaService *RuletaService) prepararNuevaRonda(ruletaActual *RuletaEnJuego) {
	ruletaActual.Mutex.Lock()
	defer ruletaActual.Mutex.Unlock()

	ruletaActual.Jugadas = nil
	ruletaActual.TimerActivo = false
	ruletaActual.Estado = EstadoEsperandoApuestas
}

func (ruletaService *RuletaService) anunciarCierreApuestas(ruletaActual *RuletaEnJuego) {
	for _, conexion := range ruletaActual.ConexionesWS {
		if conexion != nil {
			conexion.WriteJSON(map[string]string{
				"message": "¡No va más!",
			})
		}
	}
}

func (ruletaService *RuletaService) evaluarJugadas(jugadas []JugadaConUsuario, numeroGanador NumeroRuleta) map[uint][]dto.ResultadoIndividualDTO {
	resultados := make(map[uint][]dto.ResultadoIndividualDTO)

	for _, jugada := range jugadas {
		usuarioID := jugada.UsuarioID
		usuario, err := ruletaService.usuarioRepository.ObtenerPorID(usuarioID)
		if err != nil || usuario == nil {
			continue
		}

		apuesta := jugada.Apuesta
		multiplicador := calcularMultiplicador(apuesta, numeroGanador)

		resultado := ResultadoRuleta{
			MontoApostado: apuesta.Monto,
			Ganancia:      apuesta.Monto * multiplicador,
			NumeroGanador: numeroGanador,
		}

		if err := ruletaService.procesarResultado(usuario, resultado); err != nil {
			continue
		}
		if err := ruletaService.registrarJugada(usuarioID, apuesta, resultado); err != nil {
			continue
		}

		resultados[usuarioID] = append(resultados[usuarioID], dto.ResultadoIndividualDTO{
			TipoApuesta:   apuesta.TipoApuesta,
			MontoApostado: apuesta.Monto,
			Ganancia:      resultado.Ganancia,
			Detalles:      extraerDetalles(apuesta),
		})
	}

	return resultados
}

func (ruletaService *RuletaService) enviarResultados(ruletaActual *RuletaEnJuego, resultados map[uint][]dto.ResultadoIndividualDTO, numeroGanador NumeroRuleta) {
	for userID, resumen := range resultados {
		conexion := ruletaActual.ConexionesWS[userID]
		if conexion == nil {
			continue
		}
		conexion.WriteJSON(dto.RuletaResultadoUsuarioDTO{
			Mensaje:       "La ruleta ha girado",
			NumeroGanador: numeroGanador.Valor,
			Color:         numeroGanador.Color,
			Resultados:    resumen,
		})
	}
}
