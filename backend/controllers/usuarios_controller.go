// controllers/usuarios_controller.go
package controllers

import (
	"casino/db"
	"casino/dto"
	"casino/repositories"
	"casino/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

var usuarioRepo = repositories.NewUsuarioRepository(db.DB)

// var usuarioService = services.NewUsuarioService()

func CrearUsuario(c *gin.Context) {
	var input dto.CrearUsuarioDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos faltantes"})
		return
	}

	usuario, err := services.NewUsuarioService().CrearUsuario(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

// Devuelve la lista de todos los usuarios (un get all users)
func ObtenerUsuarios(c *gin.Context) {
	usuarios, err := usuarioRepo.ObtenerTodos()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	c.JSON(200, usuarios)
}
