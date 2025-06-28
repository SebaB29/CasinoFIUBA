package blackjack

import (
	"github.com/gorilla/websocket"
	"casino/models"
	"sync"
)

type BlackjackEnJuego struct {
	Mesas     map[uint]*models.MesaBlackjack
	ConexionesWS map[uint]*websocket.Conn
	Mutex        sync.Mutex
}

// NuevoBlackjackEnJuego crea la estructura global del juego en memoria
func NuevoBlackjackEnJuego() *BlackjackEnJuego {
	return &BlackjackEnJuego{
		Mesas:     make(map[uint]*models.MesaBlackjack),
		ConexionesWS: make(map[uint]*websocket.Conn),
	}
}