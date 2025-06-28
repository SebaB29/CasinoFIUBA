package protocols

import "encoding/json"

// BlackjackWSMessage es el tipo de mensaje que recibimos por WS para el blackjack
type BlackjackWSMessage struct {
	Action string          `json:"action"`
	IDMesa uint            `json:"id_mesa"`
	Data   json.RawMessage `json:"data"`
}

// ParseBlackjackWSMessage parsea los mensajes WS que llegan al blackjack
func ParseBlackjackWSMessage(msg []byte) (*BlackjackWSMessage, error) {
	var mensaje BlackjackWSMessage
	err := json.Unmarshal(msg, &mensaje)
	if err != nil {
		return nil, err
	}
	return &mensaje, nil
}
