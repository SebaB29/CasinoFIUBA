package routes

import (
	juegos "casino/controllers/juegos"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func VasosRoutes(rg *gin.RouterGroup) {
	ruta := rg.Group("/vasos")
	ruta.Use(middleware.JWTAuthMiddleware())

	ruta.POST("/nueva", juegos.CrearPartidaVasos)
	ruta.POST("/jugar", juegos.JugarVasos)
	ruta.GET("/:id", juegos.ObtenerResultadoVasos)
}
