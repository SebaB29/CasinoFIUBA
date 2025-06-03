package routes

import (
	"casino/controllers"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func JuegosRoutes(rg *gin.RouterGroup) {
	plinkoController := controllers.NewPlinkoController()
	auth := rg.Group("/juegos")
	auth.Use(middleware.JWTAuthMiddleware())

	auth.POST("/plinko", plinkoController.Jugar)
}
