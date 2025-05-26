package controllers

import (
	"casino/models"
	"github.com/gin-gonic/gin"
)

// Maneja la creación de una nueva apuesta
func CrearApuesta(c *gin.Context) {
	var a models.Apuesta
	
	// BindJSON enlaza el cuerpo de la solicitud JSON a la estructura Apuesta (de lo contrario: error 400)
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	// TODO: lógica de validacion y guardado en la base de datos
	c.JSON(201, a)
}
