package controllers

import (
	"casino/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlinkoController struct{}

func NewPlinkoController() *PlinkoController {
	return &PlinkoController{}
}

func (ctrl *PlinkoController) Jugar(c *gin.Context) {
	const monto = 100.0 // Monto fijo para test

	resultado := services.JugarPlinko(monto)

	c.JSON(http.StatusOK, gin.H{
		"mensaje":           "Jugada realizada",
		"monto_apostado":    monto,
		"posicion_final":    resultado.PosicionFinal,
		"multiplicador":     resultado.Multiplicador,
		"ganancia_obtenida": resultado.Ganancia,
	})
}
