package juegos

import (
	"casino/services/juegos/blackjack"
	protocolo "casino/websocket/protocols"
	"github.com/gorilla/websocket"
)

type BlackjackSocketHandler struct {
	conexion      *websocket.Conn
	userID        uint
	blackjackSvc  *blackjack.BlackjackService
	hub           *blackjack.BlackjackHub
}

func NewBlackjackSocketHandler(
	conexion *websocket.Conn,
	userID uint,
	service *blackjack.BlackjackService,
	hub *blackjack.BlackjackHub,
) *BlackjackSocketHandler {
	return &BlackjackSocketHandler{
		conexion:     conexion,
		userID:       userID,
		blackjackSvc: service,
		hub:          hub,
	}
}

// Manejar lee los mensajes WS y dispara la lógica correspondiente
func (handler *BlackjackSocketHandler) Manejar() {
	defer handler.conexion.Close()

	for {
		_, msg, err := handler.conexion.ReadMessage()
		if err != nil {
			break
		}

		request, err := protocolo.ParseBlackjackWSMessage(msg)
		if err != nil {
			handler.conexion.WriteJSON(map[string]string{"error": "Mensaje inválido"})
			continue
		}

		switch request.Action {
		case "hit":
			handler.blackjackSvc.Hit(handler.userID)
		case "stand":
			handler.blackjackSvc.Stand(handler.userID)
		case "doblar":
			handler.blackjackSvc.Doblar(handler.userID)
		case "rendirse":
			handler.blackjackSvc.Rendirse(handler.userID)
		case "split":
			handler.blackjackSvc.Split(handler.userID)
		case "seguro":
			handler.blackjackSvc.Seguro(handler.userID)
		default:
			handler.conexion.WriteJSON(map[string]string{"error": "Acción no reconocida"})
		}
	}
}
