package routes

import (
	"casino/controllers"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func UsuarioRoutes(rg *gin.RouterGroup) {
	usuarioController := controllers.NewUsuarioController()
	usuarios := rg.Group("/usuarios")

	// RUTAS PÚBLICAS
	usuarios.POST("/registro", usuarioController.CrearUsuario)
	usuarios.POST("/login", usuarioController.LoginUsuario)

	// RUTAS ADMIN
	usuarios.GET("/", middleware.JWTAuthMiddleware("admin"), usuarioController.ObtenerTodosLosUsuarios)
	usuarios.GET("/:id", middleware.JWTAuthMiddleware("admin"), usuarioController.ObtenerUsuarioPorID)
}
