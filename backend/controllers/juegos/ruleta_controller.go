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
	var jugada dto.RuletaRequestDTO
	if err := ctx.ShouldBindJSON(&jugada); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido"})
		return
	}

	userID := ctx.GetUint("userID")

	if err := ctrl.ruletaService.Jugar(userID, jugada); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"mensaje": "Apuesta recibida y en proceso. Esperá la ruleta..."})
}
