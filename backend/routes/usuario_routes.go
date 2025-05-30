package routes

import (
	"casino/controllers"

	"github.com/gin-gonic/gin"
)

func UsuarioRoutes(rg *gin.RouterGroup) {
	usuarios := rg.Group("/usuarios")

	// Crea instancia del controlador correctamente
	usuarioController := controllers.NewUsuarioController()

	// RUTAS PÚBLICAS
	usuarios.POST("/registro", usuarioController.CrearUsuario)
	usuarios.POST("/login", usuarioController.LoginUsuario)

	// NUEVAS RUTAS GET
	usuarios.GET("/", usuarioController.ObtenerTodosLosUsuarios)
	usuarios.GET("/:id", usuarioController.ObtenerUsuarioPorID)

	// RUTAS PROTEGIDAS (mas adelante)
	// usuarios.GET("/perfil", middleware.JWTAuthMiddleware(), usuarioController.PerfilUsuario)
}
