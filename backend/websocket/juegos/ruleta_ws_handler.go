package websocket

import (
	dto "casino/dto/juegos"
	"casino/services/ruleta"
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
			handler.conexion.WriteJSON(map[string]string{"error": "Mensaje inválido"})
			continue
		}

		switch request.Action {
		case "apostar":
			handler.procesarApuesta(request.Data)
		case "retirarse":
			handler.desconectar()
			return
		default:
			handler.responderError("Acción no reconocida")
		}
	}
}

func (handler *RuletaSocketHandler) procesarApuesta(data json.RawMessage) {
	var datos dto.RuletaRequestDTO
	if err := json.Unmarshal(data, &datos); err != nil {
		handler.conexion.WriteJSON(map[string]string{"error": "Datos de apuesta inválidos"})
		return
	}

	canal, err := handler.ruletaService.Jugar(handler.userID, datos)
	if err != nil {
		handler.conexion.WriteJSON(map[string]string{"error": err.Error()})
		return
	}

	go func() {
		resultado := <-canal
		handler.conexion.WriteJSON(dto.RuletaResponseDTO{
			Mensaje:       "La ruleta ha girado",
			NumeroGanador: resultado.NumeroGanador.Valor,
			ColorGanador:  resultado.NumeroGanador.Color,
			MontoApostado: resultado.MontoApostado,
			Ganancia:      resultado.Ganancia,
		})
	}()
}

func (handler *RuletaSocketHandler) desconectar() {
	handler.conexion.WriteJSON(map[string]string{"message": "Desconectando..."})
}

func (handler *RuletaSocketHandler) responderError(msg string) {
	handler.conexion.WriteJSON(map[string]string{"error": msg})
}
