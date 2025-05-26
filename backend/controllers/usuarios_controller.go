// controllers/usuarios_controller.go
package controllers

import (
	"casino/db"
	"casino/dto"
	"casino/errores"
	"casino/repositories"
	"casino/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

var usuarioRepo = repositories.NewUsuarioRepository(db.DB)

// var usuarioService = services.NewUsuarioService() DESPUES HAY QUE REVISAR Y CORREGIRLO POR #1

func CrearUsuario(c *gin.Context) {
	var input dto.CrearUsuarioDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos faltantes"})
		return
	}

	// #1: ESTO NO ES LO IDEAL PORQUE SE CREA UNA INSTANCIA DEL SERVICE CADA VEZ QUE SE USA, HABRIA QUE HACER UN SETUP
	usuario, err := services.NewUsuarioService().CrearUsuario(input)
	if err != nil {
		switch err {
		case errores.ErrMenorDeEdad, errores.ErrEmailYaExiste:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

func LoginUsuario(c *gin.Context) {
	var input dto.LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos faltantes"})
		return
	}

	usuario, err := services.NewUsuarioService().Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Devuelve datos del usuario sin contraseña
	c.JSON(http.StatusOK, gin.H{
		"id":      usuario.ID,
		"nombre":  usuario.Nombre,
		"email":   usuario.Email,
		"mensaje": "Inicio de sesión exitoso",
	})
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
