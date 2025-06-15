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

	resultadoChannel, err := ctrl.ruletaService.Jugar(userID, jugada)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado := <-resultadoChannel

	respuesta := dto.RuletaResponseDTO{
		Mensaje:       "La ruleta ha girado",
		NumeroGanador: resultado.NumeroGanador.Valor,
		ColorGanador:  resultado.NumeroGanador.Color,
		MontoApostado: resultado.MontoApostado,
		Ganancia:      resultado.Ganancia,
	}

	ctx.JSON(http.StatusOK, respuesta)
}
