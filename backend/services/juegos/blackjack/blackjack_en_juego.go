package blackjack

import (
	"github.com/gorilla/websocket"
	"casino/models"
	"sync"
)

type BlackjackEnJuego struct {
	Partidas     map[uint]*models.PartidaBlackjack
	ConexionesWS map[uint]*websocket.Conn
	Mutex        sync.Mutex
}

// NuevoBlackjackEnJuego crea la estructura global del juego en memoria
func NuevoBlackjackEnJuego() *BlackjackEnJuego {
	return &BlackjackEnJuego{
		Partidas:     make(map[uint]*models.PartidaBlackjack),
		ConexionesWS: make(map[uint]*websocket.Conn),
	}
}
