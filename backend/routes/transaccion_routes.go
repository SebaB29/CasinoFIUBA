package routes

import (
	"casino/controllers"
	"casino/middleware"

	"github.com/gin-gonic/gin"
)

func TransaccionRoutes(rg *gin.RouterGroup) {
	transaccionCtrl := controllers.NewTransaccionController()
	auth := rg.Group("/transacciones")
	auth.Use(middleware.JWTAuthMiddleware())

	auth.POST("/depositar", transaccionCtrl.Depositar)
	auth.POST("/extraer", transaccionCtrl.Extraer)
}
