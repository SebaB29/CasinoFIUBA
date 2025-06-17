package routes

import (
	controllers "casino/controllers/juegos"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func BuscaminasRoutes(rg *gin.RouterGroup) {
	ruta := rg.Group("/buscaminas")
	ruta.Use(middleware.JWTAuthMiddleware())

	ruta.POST("/nueva", controllers.CrearPartidaBuscaminas)
	ruta.POST("/abrir", controllers.AbrirCeldaBuscaminas)
	ruta.POST("/retirarse", controllers.RetirarseBuscaminas)
	ruta.GET("/debug/:id", controllers.VerMinasDebug)
}
