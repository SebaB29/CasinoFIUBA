
package controllers

import (
	"casino/juegos/buscaminas"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AbrirCeldaRequest struct {
	Fila int `json:"fila" binding:"required,min=0,max=4"`
	Col  int `json:"col" binding:"required,min=0,max=4"`
}

type NuevaPartidaRequest struct {
	Minas   int     `json:"minas" binding:"required,min=1,max=24"`
	Apuesta float64 `json:"apuesta" binding:"required,gt=0"`
}

func CrearPartidaBuscaminas(c *gin.Context) {
	var req NuevaPartidaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partida, err := buscaminas.CrearPartida(5, 5, req.Minas, req.Apuesta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"estado":          partida.Estado,
		"minas":           partida.CantidadMinas,
		"monto_apostado":  partida.MontoApostado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})

	buscaminas.SetPartidaActual(partida) // Solo para pruebas, se elimina en la versión final

}

func AbrirCeldaBuscaminas(c *gin.Context) {
	var req AbrirCeldaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partida := buscaminas.ObtenerPartidaActual()
	if partida == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay partida activa"})
		return
	}

	err := partida.AbrirCelda(req.Fila, req.Col)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"estado":          partida.Estado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})
}

func RetirarseBuscaminas(c *gin.Context) {
	partida := buscaminas.ObtenerPartidaActual()
	if partida == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay partida activa"})
		return
	}

	premio, err := partida.Retirarse()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Te retiraste exitosamente",
		"estado":  partida.Estado,
		"premio":  premio,
	})
}

// ocultarMinas evita mostrar dónde están las minas
func ocultarMinas(t *buscaminas.Tablero) [][]map[string]interface{} {
	resultado := make([][]map[string]interface{}, t.Filas)
	for i := range t.Celdas {
		resultado[i] = make([]map[string]interface{}, t.Columnas)
		for j, celda := range t.Celdas[i] {
			resultado[i][j] = map[string]interface{}{
				"abierta":        celda.Abierta,
				"marcada":        celda.Marcada,
				"minas_cercanas": celda.MinasVecinas,
			}
		}
	}
	return resultado
}
