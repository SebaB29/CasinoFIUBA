package ruleta

import (
	dto "casino/dto/juegos"
	"sync"

	"github.com/gorilla/websocket"
)

type RuletaManager struct {
	JuegoActual RuletaEnJuego
}

type RuletaEnJuego struct {
	Jugadas      []JugadaConUsuario
	Mutex        sync.Mutex
	TimerActivo  bool
	ConexionesWS map[uint]*websocket.Conn // NUEVO
}

type JugadaConUsuario struct {
	UsuarioID uint
	Apuesta   dto.RuletaRequestDTO
	Resultado chan ResultadoRuleta // Canal para enviar el resultado
	Conexion  *websocket.Conn
}

type ResultadoRuleta struct {
	MontoApostado float64
	Ganancia      float64
	NumeroGanador NumeroRuleta
}
