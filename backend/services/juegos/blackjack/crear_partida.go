package blackjack

import (
	"casino/models"
	repo "casino/repositories/juegos"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (s *BlackjackService) CrearPartida(userID uint, apuesta float64) (gin.H, error) {
	mazo := MezclarMazo(CantidadMazos)

	// Reparto inicial
	carta1, mazo := TomarCarta(mazo)
	carta2, mazo := TomarCarta(mazo)
	cartaB1, mazo := TomarCarta(mazo)
	cartaB2, mazo := TomarCarta(mazo)

	cartasJugador := []string{carta1, carta2}
	cartasBanca := []string{cartaB1, cartaB2}

	estado := models.EnCurso

	if TieneBlackjack(cartasJugador) {
		if TieneBlackjack(cartasBanca) {
			estado = models.Empatada
		} else {
			estado = models.Ganada
		}
	} else if TieneBlackjack(cartasBanca) {
		estado = models.Perdida
	}

	partida := &models.PartidaBlackjack{
		UserID:        userID,
		Apuesta:       apuesta,
		CartasJugador: CartasToString(cartasJugador),
		CartasBanca:   CartasToString(cartasBanca),
		Mazo:          CartasToString(mazo),
		Estado:        estado,
	}

	if err := repo.CrearPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo crear la partida: %w", err)
	}

	cartasBancaVisible := []string{cartaB1}
	if estado != models.EnCurso {
		cartasBancaVisible = cartasBanca
	}

	return gin.H{
		"id":             partida.ID,
		"user_id":        partida.UserID,
		"apuesta":        partida.Apuesta,
		"cartas_jugador": cartasJugador,
		"cartas_banca":   cartasBancaVisible,
		"estado":         partida.Estado,
	}, nil
}
