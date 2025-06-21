package blackjack

import (
	juegos "casino/repositories/juegos"
	"fmt"

	"github.com/gin-gonic/gin"
)

type BlackjackService struct{}

func NewBlackjackService() *BlackjackService {
	return &BlackjackService{}
}

func (s *BlackjackService) ObtenerEstado(id uint) (gin.H, error) {
	partida, err := juegos.ObtenerPartidaPorID(id)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}

	return gin.H{
		"id":             partida.ID,
		"cartas_mano_1":  StringToCartas(partida.CartasJugador),
		"cartas_mano_2":  StringToCartas(partida.CartasJugadorSplit),
		"cartas_banca":   StringToCartas(partida.CartasBanca),
		"estado":         partida.Estado,
	}, nil
}
