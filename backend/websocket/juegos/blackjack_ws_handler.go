package juegos

import (
	"casino/services/juegos/blackjack"
	protocolo "casino/websocket/protocols"
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
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

        var resp gin.H
        var actionErr error

        switch request.Action {
        case "hit":
            resp, actionErr = handler.blackjackSvc.Hit(request.IDPartida)
        case "stand":
            resp, actionErr = handler.blackjackSvc.Stand(request.IDPartida)
        case "doblar":
            resp, actionErr = handler.blackjackSvc.Doblar(request.IDPartida)
        case "rendirse":
            resp, actionErr = handler.blackjackSvc.Rendirse(request.IDPartida)
        case "split":
            resp, actionErr = handler.blackjackSvc.Split(request.IDPartida)
        case "seguro":
            resp, actionErr = handler.blackjackSvc.Seguro(request.IDPartida)
        default:
            handler.conexion.WriteJSON(map[string]string{"error": "Acción no reconocida"})
            continue
        }

        if actionErr != nil {
            handler.conexion.WriteJSON(map[string]string{"error": actionErr.Error()})
            continue
        }

        // Mandar estado al cliente que envió la acción
        handler.conexion.WriteJSON(resp)

        // Además, enviar el estado a todos los demás que corresponda
        handler.hub.BroadcastEstado(handler.userID, resp)
    }
}

