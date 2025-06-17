package controllers

import (
	dto "casino/dto/juegos"
	plinko "casino/services/juegos/plinko"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlinkoController struct {
	service *plinko.PlinkoService
}

func NewPlinkoController() *PlinkoController {
	return &PlinkoController{
		service: plinko.NewPlinkoService(),
	}
}

func (ctrl *PlinkoController) Jugar(c *gin.Context) {
	var input dto.PlinkoRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos o monto faltante"})
		return
	}

	userID := c.GetUint("userID")
	resultado, err := ctrl.service.Jugar(userID, input.Monto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultado)
}
