package controllers

import (
	dto "casino/dto/juegos"
	ruleta "casino/services/ruleta"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RuletaController struct {
	ruletaService *ruleta.RuletaService
}

func NewRuletaController() *RuletaController {
	return &RuletaController{
		ruletaService: ruleta.NewRuletaService(),
	}
}

func (ctrl *RuletaController) Jugar(ctx *gin.Context) {
	var request dto.RuletaRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de apuesta inválidos"})
		return
	}

	userID := ctx.GetUint("userID")
	resultado, err := ctrl.ruletaService.Jugar(userID, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}
