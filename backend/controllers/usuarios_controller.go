// controllers/usuarios_controller.go
package controllers

import (
	"casino/dto"
	"casino/errores"
	"casino/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UsuarioController struct {
	service services.UsuarioServiceInterface
}

func NewUsuarioController() *UsuarioController {
	service := services.NewUsuarioService()
	return &UsuarioController{service: service}
}

func (ctrl *UsuarioController) CrearUsuario(c *gin.Context) {
	var input dto.CrearUsuarioDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos faltantes"})
		return
	}

	usuario, err := ctrl.service.CrearUsuario(input)
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

func (ctrl *UsuarioController) LoginUsuario(c *gin.Context) {
	var input dto.LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos faltantes"})
		return
	}

	usuario, token, err := ctrl.service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      usuario.ID,
		"nombre":  usuario.Nombre,
		"email":   usuario.Email,
		"token":   token,
		"mensaje": "Inicio de sesión exitoso",
	})
}

// Devuelve la lista de todos los usuarios (un get all users)
// func ObtenerUsuarios(c *gin.Context) {
// 	usuarios, err := usuarioRepo.ObtenerTodos()
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "Error al obtener usuarios"})
// 		return
// 	}

// 	c.JSON(200, usuarios)
// }
