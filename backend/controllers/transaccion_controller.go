package controllers

import (
	"casino/dto"
	"casino/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransaccionController struct {
	service services.TransaccionServiceInterface
}

func NewTransaccionController() *TransaccionController {
	return &TransaccionController{
		service: services.NewTransaccionService(),
	}
}

func (ctrl *TransaccionController) Depositar(c *gin.Context) {
	var input dto.TransaccionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userID := c.GetUint("userID")
	err := ctrl.service.Depositar(userID, input.Monto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Depósito exitoso"})
}

func (ctrl *TransaccionController) Extraer(c *gin.Context) {
	var input dto.TransaccionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userID := c.GetUint("userID")
	err := ctrl.service.Extraer(userID, input.Monto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Extracción exitosa"})
}
