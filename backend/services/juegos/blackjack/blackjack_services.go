package blackjack

import (
	"fmt"
	"github.com/gin-gonic/gin"
	repo "casino/repositories/juegos"
)

// BlackjackService contiene el estado en memoria y el hub para eventos
type BlackjackService struct {
	enJuego *BlackjackEnJuego
	Hub     *BlackjackHub
}

// NuevoBlackjackService crea el servicio con enJuego y el hub
func NewBlackjackService() *BlackjackService {
	enJuego := NuevoBlackjackEnJuego()
	hub := NuevoBlackjackHub(enJuego)

	// Ejecuta el hub en goroutine
	go hub.Run()

	return &BlackjackService{
		enJuego: enJuego,
		Hub:     hub,
	}
}

// CrearPartida crea una nueva partida en base de datos y en memoria
func (s *BlackjackService) CrearPartida(userID uint, apuesta float64) (gin.H, error) {
	resp, err := NuevaPartida(userID, apuesta)
	if err != nil {
		return nil, err
	}
	// Cargar partida
	id := resp["id"].(uint)
	partida, err := repo.ObtenerPartidaPorID(id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener partida: %w", err)
	}
	// Registra en memoria
	s.enJuego.Mutex.Lock()
	s.enJuego.Partidas[userID] = partida
	s.enJuego.Mutex.Unlock()

	// Mandar estado por WS
	s.Hub.BroadcastEstado(userID, resp)
	return resp, nil
}

// ObtenerEstado devuelve el estado de la partida
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
		"cartas_banca":        []string{cartasBanca[0], "???"},
		"valor_banca_visible": valorBancaVisible, 
	}, nil
}
