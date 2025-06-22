package controllers

import (
	dto "casino/dto/juegos"
	"casino/services/juegos/slot"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SlotController struct {
	service *slot.SlotService
}

func NewSlotController() *SlotController {
	return &SlotController{
		service: slot.NewSlotService(),
	}
}

func (ctrl *SlotController) Jugar(ctx *gin.Context) {
	var input dto.SlotRequestDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos o monto faltante"})
		return
	}

	userID := ctx.GetUint("userID")
	resultado, err := ctrl.service.Jugar(userID, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}
