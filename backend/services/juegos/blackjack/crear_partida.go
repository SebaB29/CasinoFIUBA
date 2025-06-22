package blackjack

import (
    "casino/models"
    repo "casino/repositories/juegos"
    "fmt"
    "github.com/gin-gonic/gin"
)

// NuevaPartida es la función que encapsula la creación de una nueva partida.
func NuevaPartida(userID uint, apuesta float64) (gin.H, error) {
    if apuesta <= 0 {
        return nil, fmt.Errorf("la apuesta debe ser mayor a cero")
    }

    mazo := MezclarMazo(CantidadMazos)

    // Reparto inicial
    carta1Jugador, mazo := TomarCarta(mazo)
    carta2Jugador, mazo := TomarCarta(mazo)
    carta1Banca, mazo := TomarCarta(mazo)
    carta2Banca, mazo := TomarCarta(mazo)

    cartasJugador := []string{carta1Jugador, carta2Jugador}
    cartasBanca := []string{carta1Banca, carta2Banca}

    estado := models.EnCurso
    blackjackNatural := false

    tieneJugadorBlackjack := TieneBlackjack(cartasJugador)
    tieneBancaBlackjack := TieneBlackjack(cartasBanca)

    switch {
    case tieneJugadorBlackjack && tieneBancaBlackjack:
        estado = models.Empatada
        blackjackNatural = true
    case tieneJugadorBlackjack:
        estado = models.Ganada
        blackjackNatural = true
    case tieneBancaBlackjack:
        estado = models.Perdida
        blackjackNatural = true
    }

    partida := &models.PartidaBlackjack{
        UserID:           userID,
        Apuesta:          apuesta,
        CartasJugador:    CartasToString(cartasJugador),
        CartasBanca:      CartasToString(cartasBanca),
        Mazo:             CartasToString(mazo),
        Estado:           estado,
        ManoActual:       1,
        BlackjackNatural: blackjackNatural,
        Doblar:           false,
        Seguro:           0,
    }

    if err := repo.CrearPartida(partida); err != nil {
        return nil, fmt.Errorf("no se pudo crear la partida: %w", err)
    }

    // Cartas de la banca visibles solo si ya terminó
    var cartasBancaParaCliente []string
    valorBancaParaCliente := 0
    if partida.Estado == models.EnCurso {
        cartasBancaParaCliente = []string{carta1Banca, "???"}
        valorBancaParaCliente = CalcularValor([]string{carta1Banca})
    } else {
        cartasBancaParaCliente = cartasBanca
        valorBancaParaCliente = CalcularValor(cartasBanca)
    }

    return gin.H{
        "id":                  partida.ID,
        "user_id":             partida.UserID,
        "apuesta_principal":   partida.Apuesta,
        "cartas_jugador":      cartasJugador,
        "valor_jugador":       CalcularValor(cartasJugador),
        "cartas_banca":        cartasBancaParaCliente,
        "valor_banca_visible": valorBancaParaCliente,
        "estado":              partida.Estado,
        "blackjack_natural":   partida.BlackjackNatural,
        "mensaje":             "Partida creada",
    }, nil
}
