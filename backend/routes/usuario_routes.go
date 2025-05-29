package routes

import (
	"casino/controllers"

	"github.com/gin-gonic/gin"
)

func RegistroUsuarioRoutes(rg *gin.RouterGroup) {
	usuarios := rg.Group("/usuarios")

	// RUTAS PÚBLICAS
	usuarios.POST("/registro", usuarioController.CrearUsuario)
	usuarios.POST("/login", usuarioController.LoginUsuario)

	// NUEVAS RUTAS GET
	usuarios.GET("/", usuarioController.ObtenerTodosLosUsuarios)
	usuarios.GET("/:id", usuarioController.ObtenerUsuarioPorID)

	// RUTAS PROTEGIDAS (futuro)
	// usuarios.GET("/perfil", middleware.JWTAuthMiddleware(), usuarioController.PerfilUsuario)
}
