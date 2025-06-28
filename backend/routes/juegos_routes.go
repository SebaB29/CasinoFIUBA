package routes

import (
	controllers "casino/controllers/juegos"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func JuegosRoutes(rg *gin.RouterGroup) {
	plinkoController := controllers.NewPlinkoController()
	ruletaController := controllers.NewRuletaController()
	slotController := controllers.NewSlotController()
	auth := rg.Group("/juegos")
	auth.Use(middleware.JWTAuthMiddleware())

	auth.POST("/slot", slotController.Jugar)
	auth.POST("/plinko", plinkoController.Jugar)
	auth.GET("/ruleta", ruletaController.JugarWS)

}
