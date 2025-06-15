package ruleta

import (
	dto "casino/dto/juegos"
	"sync"
)

type RuletaManager struct {
	JuegoActual RuletaEnJuego
}

type RuletaEnJuego struct {
	Jugadas     []JugadaConUsuario
	Mutex       sync.Mutex
	TimerActivo bool
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
