package routes

import (
	"casino/controllers"

	"github.com/gin-gonic/gin"
)

func RegistroUsuarioRoutes(rg *gin.RouterGroup) {
	usuarios := rg.Group("/usuarios")

	usuarios.POST("/registro", controllers.CrearUsuario)
	usuarios.POST("/login", controllers.LoginUsuario)
	usuarios.GET("/", controllers.ObtenerUsuarios)
}
