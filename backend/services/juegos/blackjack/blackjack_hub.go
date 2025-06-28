package blackjack

import (
	dto "casino/dto/juegos"
	"github.com/gin-gonic/gin"
)

// EventoBlackjack representa un evento que se envía por WS a un jugador
type EventoBlackjack struct {
	UserID uint
	Estado gin.H
}

// BlackjackHub mantiene un canal para eventos de juego
type BlackjackHub struct {
	eventos chan EventoBlackjack
	juego   *BlackjackEnJuego
}

// NuevoBlackjackHub crea una nueva instancia del hub
func NuevoBlackjackHub(juego *BlackjackEnJuego) *BlackjackHub {
	return &BlackjackHub{
		eventos: make(chan EventoBlackjack, 100),
		juego:   juego,
	}
}

// Run procesa los eventos y los envía por WS
func (hub *BlackjackHub) Run() {
	for evento := range hub.eventos {
		hub.juego.Mutex.Lock()
		conn, ok := hub.juego.ConexionesWS[evento.UserID]
		hub.juego.Mutex.Unlock()

		if ok && conn != nil {
			conn.WriteJSON(evento.Estado)
		}
	}
}

// BroadcastEstado encola un evento para enviar estado a un jugador
func (hub *BlackjackHub) BroadcastEstado(userID uint, estado dto.BlackjackEstadoDTO) {
	hub.eventos <- EventoBlackjack{
		UserID: userID,
		Estado: gin.H{"estado": estado},
	}
}
