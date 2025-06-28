package websocket

import "encoding/json"

type mensajeWS struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

func ParseWSMessage(msg []byte) (*mensajeWS, error) {
	var mensaje mensajeWS
	err := json.Unmarshal(msg, &mensaje)
	if err != nil {
		return nil, err
	}
	return &mensaje, nil
}
