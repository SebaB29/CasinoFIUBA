package controllers

import (
	"casino/db"
	dto "casino/dto/juegos"
	"casino/juegos/vasos"
	"casino/models"
	repo "casino/repositories/juegos"
	"github.com/gin-gonic/gin"
	"net/http"
)

// POST /vasos/nueva
func CrearPartidaVasos(c *gin.Context) {
	var input dto.CrearPartidaVasosDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	logica := vasos.NuevaPartidaVasos(input.Apuesta)

	model := &models.PartidaVasos{
		UserID:           userID,
		PosicionCorrecta: logica.PosicionCorrecta,
		Apuesta:          logica.Apuesta,
		Estado:           string(logica.Estado),
		FechaCreacion:    logica.CreadoEn,
	}

	if err := repo.CrearPartidaVasos(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la partida", "detalle": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

// POST /vasos/jugar
func JugarVasos(c *gin.Context) {
	var input dto.RealizarJugadaVasosDTO
	if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
        "raw": input,
    })
    return
}


	partida, err := repo.ObtenerPartidaVasosPorID(input.IDPartida)
	if err != nil || partida == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	if partida.Estado != string(vasos.EnCurso) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La partida ya fue jugada"})
		return
	}

	acierto := *input.Eleccion == partida.PosicionCorrecta
	if acierto {
		partida.Estado = string(vasos.Ganada)
	} else {
		partida.Estado = string(vasos.Perdida)
	}
	partida.PosicionElegida = input.Eleccion

	if err := repo.ActualizarPartidaVasos(partida); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la partida"})
		return
	}

	resultado := "perdiste"
	if acierto {
		resultado = "ganaste"
	}

	c.JSON(http.StatusOK, gin.H{
		"resultado":         resultado,
		"posicion_correcta": partida.PosicionCorrecta,
		"posicion_elegida":  *partida.PosicionElegida,
		"estado":            partida.Estado,
	})
}

// GET /vasos/:id
func ObtenerResultadoVasos(c *gin.Context) {
	id := c.Param("id")

	var partida models.PartidaVasos
	if err := db.DB.First(&partida, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	c.JSON(http.StatusOK, partida)
}
