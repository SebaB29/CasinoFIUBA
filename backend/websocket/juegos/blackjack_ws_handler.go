package websocket

import (
	juegos "casino/dto/juegos"
	"casino/services/juegos/blackjack"
	protocolo "casino/websocket/protocols"

	"github.com/gorilla/websocket"
)

type BlackjackSocketHandler struct {
	conexion     *websocket.Conn
	userID       uint
	blackjackSvc *blackjack.BlackjackService
	hub          *blackjack.BlackjackHub
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

// Manejar procesa los mensajes del WebSocket para el juego de Blackjack
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

		var resp juegos.BlackjackEstadoDTO
		var actionErr error

		idMesa := request.IDMesa

		switch request.Action {
		case "hit":
			resp, actionErr = handler.blackjackSvc.Hit(idMesa, handler.userID)
		case "stand":
			resp, actionErr = handler.blackjackSvc.Stand(idMesa, handler.userID)
		case "doblar":
			resp, actionErr = handler.blackjackSvc.Doblar(idMesa, handler.userID)
		case "rendirse":
			resp, actionErr = handler.blackjackSvc.Rendirse(idMesa, handler.userID)
		case "split":
			resp, actionErr = handler.blackjackSvc.Split(idMesa, handler.userID)
		case "seguro":
			resp, actionErr = handler.blackjackSvc.Seguro(idMesa, handler.userID)
		default:
			handler.conexion.WriteJSON(map[string]string{"error": "Acción no reconocida"})
			continue
		}

		if actionErr != nil {
			handler.conexion.WriteJSON(map[string]string{"error": actionErr.Error()})
			continue
		}

		// Enviar al jugador actual
		handler.conexion.WriteJSON(resp)

		// Estado del jugador que hizo la jugada
		estado, err := handler.blackjackSvc.ObtenerEstadoMesa(idMesa, handler.userID)
		if err == nil {
			handler.hub.BroadcastEstado(handler.userID, estado)
		}

		// Broadcast al jugador al que le toca el turno
		mesa, err := handler.blackjackSvc.ObtenerMesa(idMesa)
		if err == nil && mesa.Estado == "en_curso" && mesa.JugadorActual < len(mesa.ManosJugadores) {
			nuevoJugadorID := mesa.ManosJugadores[mesa.JugadorActual].UserID

			if nuevoJugadorID != handler.userID {
				estadoNuevo := handler.blackjackSvc.EstadoParaJugador(mesa, nuevoJugadorID)
				handler.hub.BroadcastEstado(nuevoJugadorID, estadoNuevo)
			}
		}
	}
}
