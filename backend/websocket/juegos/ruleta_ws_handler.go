package websocket

import (
	dto "casino/dto/juegos"
	ruleta "casino/services/juegos/ruleta"
	protocolo "casino/websocket/protocols"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type RuletaSocketHandler struct {
	conexion      *websocket.Conn
	userID        uint
	ruletaService *ruleta.RuletaService
}

func NewRuletaSocketHandler(conexion *websocket.Conn, userID uint, servicio *ruleta.RuletaService) *RuletaSocketHandler {
	return &RuletaSocketHandler{
		conexion:      conexion,
		userID:        userID,
		ruletaService: servicio,
	}
}

func (handler *RuletaSocketHandler) Manejar() {
	defer handler.conexion.Close()

	for {
		_, msg, err := handler.conexion.ReadMessage()
		if err != nil {
			break
		}

		request, err := protocolo.ParseWSMessage(msg)
		if err != nil {
			handler.responderError("Mensaje inválido")
			continue
		}

		switch request.Action {
		case "apostar":
			handler.procesarApuestas(request.Data)
		case "retirarse":
			handler.desconectar()
			return
		default:
			handler.responderError("Acción no reconocida")
		}
	}
}

func (handler *RuletaSocketHandler) procesarApuestas(data json.RawMessage) {
	var req struct {
		Apuestas []dto.RuletaRequestDTO `json:"apuestas"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handler.responderError("Datos de apuestas inválidos")
		return
	}

	if len(req.Apuestas) == 0 {
		handler.responderError("No se enviaron apuestas")
		return
	}

	for _, apuesta := range req.Apuestas {
		if err := handler.ruletaService.Jugar(handler.userID, apuesta, handler.conexion); err != nil {
			handler.responderError("Error en una de las apuestas: " + err.Error())
			return
		}
	}
}

func (handler *RuletaSocketHandler) desconectar() {
	handler.conexion.WriteJSON(map[string]string{"message": "Desconectando..."})
}

func (handler *RuletaSocketHandler) responderError(msg string) {
	handler.conexion.WriteJSON(map[string]string{"error": msg})
}
