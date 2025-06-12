package controllers

import (
	dto "casino/dto/juegos"
	services "casino/services/ruleta"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RuletaController struct {
	service *services.RuletaService
}

func NewRuletaController() *RuletaController {
	return &RuletaController{
		service: services.NewRuletaService(),
	}
}

func (ctrl *RuletaController) Jugar(ctx *gin.Context) {
	var request dto.RuletaRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de apuesta inválidos"})
		return
	}

	userID := ctx.GetUint("userID")
	resultado, err := ctrl.service.Jugar(userID, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}
