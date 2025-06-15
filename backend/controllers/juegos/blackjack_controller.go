package controllers

import (
    "casino/db"
    dto "casino/dto/juegos"
    "casino/juegos/blackjack"
    "casino/models"
    repo "casino/repositories/juegos"
    "net/http"

    "github.com/gin-gonic/gin"
)

// POST /blackjack/nueva
func CrearPartidaBlackjack(c *gin.Context) {
    var input dto.IniciarBlackjackDTO
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetUint("userID")
    mazo := blackjack.MezclarMazo()

    // Tomamos 2 para el jugador, 2 para la banca
    carta1, mazo := blackjack.TomarCarta(mazo)
    carta2, mazo := blackjack.TomarCarta(mazo)
    cartaB1, mazo := blackjack.TomarCarta(mazo)
    cartaB2, mazo := blackjack.TomarCarta(mazo)

    partida := &models.PartidaBlackjack{
        UserID:        userID,
        Apuesta:       input.Apuesta,
        CartasJugador: blackjack.CartasToString([]string{carta1, carta2}),
        CartasBanca:   blackjack.CartasToString([]string{cartaB1, cartaB2}), // Guardamos ambas
        Mazo:          blackjack.CartasToString(mazo),
        Estado:        models.EnCurso,
    }

    if err := repo.CrearPartida(partida); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la partida"})
        return
    }

    // Solo mostramos la primera carta de la banca en la respuesta inicial
    c.JSON(http.StatusOK, gin.H{
        "id":             partida.ID,
        "user_id":        partida.UserID,
        "apuesta":        partida.Apuesta,
        "cartas_jugador": []string{carta1, carta2},
        "cartas_banca":   []string{cartaB1},
        "estado":         partida.Estado,
    })
}

// POST /blackjack/hit
func HitBlackjack(c *gin.Context) {
    var input dto.JugadaBlackjackDTO
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    partida, err := repo.ObtenerPartidaPorID(input.IDPartida)
    if err != nil || partida.Estado != models.EnCurso {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Partida no encontrada o ya finalizada"})
        return
    }

    mazo := blackjack.StringToCartas(partida.Mazo)
    cartasJugador := blackjack.StringToCartas(partida.CartasJugador)

    nuevaCarta, mazo := blackjack.TomarCarta(mazo)
    cartasJugador = append(cartasJugador, nuevaCarta)

    valor := blackjack.CalcularValor(cartasJugador)
    if valor > 21 {
        partida.Estado = models.Perdida
    }

    // Actualizar partida
    partida.CartasJugador = blackjack.CartasToString(cartasJugador)
    partida.Mazo = blackjack.CartasToString(mazo)

    if err := repo.ActualizarPartida(partida); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando partida"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "cartas_jugador": cartasJugador,
        "valor":          valor,
        "estado":         partida.Estado,
    })
}

// POST /blackjack/stand
func StandBlackjack(c *gin.Context) {
    var input dto.JugadaBlackjackDTO
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    partida, err := repo.ObtenerPartidaPorID(input.IDPartida)
    if err != nil || partida.Estado != models.EnCurso {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Partida no encontrada o ya finalizada"})
        return
    }

    cartasJugador := blackjack.StringToCartas(partida.CartasJugador)
    mazo := blackjack.StringToCartas(partida.Mazo)

    // Revelar todas las cartas de la banca guardadas desde el inicio
    cartasBanca := blackjack.StringToCartas(partida.CartasBanca)

    // Reglas de la banca: sacar hasta alcanzar 17
    for blackjack.CalcularValor(cartasBanca) < 17 && len(mazo) > 0 {
        var carta string
        carta, mazo = blackjack.TomarCarta(mazo)
        cartasBanca = append(cartasBanca, carta)
    }

    valJugador := blackjack.CalcularValor(cartasJugador)
    valBanca := blackjack.CalcularValor(cartasBanca)

    switch {
    case valJugador > 21:
        partida.Estado = models.Perdida
    case valBanca > 21:
        partida.Estado = models.Ganada
    case valJugador > valBanca:
        partida.Estado = models.Ganada
    case valJugador < valBanca:
        partida.Estado = models.Perdida
    default:
        partida.Estado = models.Empatada
    }

    // Actualizar partida
    partida.CartasBanca = blackjack.CartasToString(cartasBanca)
    partida.Mazo = blackjack.CartasToString(mazo)

    if err := repo.ActualizarPartida(partida); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando estado"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "cartas_jugador": cartasJugador,
        "cartas_banca":   cartasBanca,
        "estado":         partida.Estado,
    })
}

// GET /blackjack/estado/:id_partida
func ObtenerEstadoBlackjack(c *gin.Context) {
    idParam := c.Param("id_partida")
    var partida models.PartidaBlackjack

    if err := db.DB.First(&partida, idParam).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
        return
    }

    c.JSON(http.StatusOK, partida)
}
