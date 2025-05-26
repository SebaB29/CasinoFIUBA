package controllers

import (
	"casino/db"
	"casino/models"
	"casino/repositories"
	"github.com/gin-gonic/gin"
)

var usuarioRepo = repositories.NewUsuarioRepository(db.DB)

// Guarda un nuevo usuario en la base de datos
func CrearUsuario(c *gin.Context) {
	var u models.Usuario

	// Parsea el JSON recibido al struct Usuario
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	if err := usuarioRepo.Crear(&u); err != nil {
		c.JSON(500, gin.H{"error": "Error al guardar el usuario"})
		return
	}

	c.JSON(201, u)
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
