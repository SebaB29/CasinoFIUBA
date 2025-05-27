package routes

import (
	"casino/controllers"

	"github.com/gin-gonic/gin"
)

func UsuarioRoutes(rg *gin.RouterGroup) {
	usuarios := rg.Group("/usuarios")
	usuarioController := controllers.NewUsuarioController()

	// RUTAS PÚBLICAS
	usuarios.POST("/registro", usuarioController.CrearUsuario)
	usuarios.POST("/login", usuarioController.LoginUsuario)
	// usuarios.GET("/", usuarioController.ObtenerUsuarios)

	// RUTAS PROTEGIDAS
	// EJEMPLO: usuarios.GET("/perfil", middleware.JWTAuthMiddleware(), usuarioController.PerfilUsuario)
}
