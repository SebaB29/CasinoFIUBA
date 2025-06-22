package blackjack

import (
    "fmt"
    "github.com/gin-gonic/gin"
    repo "casino/repositories/juegos"
)

type BlackjackService struct{}

func NewBlackjackService() *BlackjackService {
    return &BlackjackService{}
}

// CrearPartida delega la lógica a la función helper en crear_partida.go
func (s *BlackjackService) CrearPartida(userID uint, apuesta float64) (gin.H, error) {
    return NuevaPartida(userID, apuesta)
}

// ObtenerEstado retorna el estado actual de una partida
func (s *BlackjackService) ObtenerEstado(idPartida uint) (gin.H, error) {
    partida, err := repo.ObtenerPartidaPorID(idPartida)
    if err != nil {
        return nil, fmt.Errorf("partida no encontrada")
    }

    cartasBanca := StringToCartas(partida.CartasBanca)
    valorBancaVisible := 0
    if len(cartasBanca) > 0 {
        valorBancaVisible = CalcularValor([]string{cartasBanca[0]})
    }

    return gin.H{
        "id":                  partida.ID,
        "user_id":             partida.UserID,
        "estado":              partida.Estado,
        "mano_actual":         partida.ManoActual,
        "cartas_jugador":      StringToCartas(partida.CartasJugador),
        "valor_jugador":       CalcularValor(StringToCartas(partida.CartasJugador)),
        "cartas_jugador_split": StringToCartas(partida.CartasJugadorSplit),
        "valor_jugador_split": CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
        "cartas_banca":        []string{cartasBanca[0], "???"}, // solo visible
        "valor_banca_visible": valorBancaVisible,
    }, nil
}
