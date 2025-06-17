package routes

import (
	controllers "casino/controllers/juegos"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func JuegosRoutes(rg *gin.RouterGroup) {
	plinkoController := controllers.NewPlinkoController()
	ruletaController := controllers.NewRuletaController()
	auth := rg.Group("/juegos")
	auth.Use(middleware.JWTAuthMiddleware())

	auth.POST("/plinko", plinkoController.Jugar)
	auth.POST("/ruleta", ruletaController.Jugar)
}
