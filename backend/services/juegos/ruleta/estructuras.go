package ruleta

import (
	dto "casino/dto/juegos"
	"sync"

	"github.com/gorilla/websocket"
)

type EstadoRuleta string

const (
	EstadoEsperandoApuestas EstadoRuleta = "esperando"
	EstadoGirando           EstadoRuleta = "girando"
)

type RuletaEnJuego struct {
	Jugadas      []JugadaConUsuario
	ConexionesWS map[uint]*websocket.Conn // Jugadores: ConexionesWS[UserID]
	TimerActivo  bool
	Estado       EstadoRuleta
	Mutex        sync.Mutex
}

type JugadaConUsuario struct {
	UsuarioID uint
	Apuesta   dto.RuletaRequestDTO
}

type ResultadoRuleta struct {
	MontoApostado float64
	Ganancia      float64
	NumeroGanador NumeroRuleta
}
