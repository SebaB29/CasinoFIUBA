package routes

import (
	"casino/controllers"

	"github.com/gin-gonic/gin"
)

func RegistroUsuarioRoutes(rg *gin.RouterGroup) {
	usuarios := rg.Group("/usuarios")

	usuarios.POST("/registro", controllers.CrearUsuario)
	usuarios.GET("/", controllers.ObtenerUsuarios)
}
