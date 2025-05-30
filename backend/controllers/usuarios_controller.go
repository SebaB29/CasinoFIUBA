// controllers/usuarios_controller.go
package controllers

import (
	"casino/dto"
	"casino/errores"
	"casino/services"
	"net/http"
	"strconv"

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
	var input dto.RegistroUsuarioRequestDTO

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
	var input dto.LoginRequestDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o campos faltantes"})
		return
	}

	usuario, err := ctrl.service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      usuario.ID,
		"nombre":  usuario.Nombre,
		"email":   usuario.Email,
		"token":   usuario.Token,
		"mensaje": "Inicio de sesión exitoso",
	})
}

// GET /usuarios
func (ctrl *UsuarioController) ObtenerTodosLosUsuarios(c *gin.Context) {
	usuarios, err := ctrl.service.ObtenerTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	c.JSON(http.StatusOK, usuarios)
}

// GET /usuarios (por id)
func (ctrl *UsuarioController) ObtenerUsuarioPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	usuario, err := ctrl.service.ObtenerPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usuario)
}
