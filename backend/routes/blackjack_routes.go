package routes

import (
	juegos "casino/controllers/juegos"
	"github.com/gin-gonic/gin"
	"casino/middleware"
)

func BlackjackRoutes(rg *gin.RouterGroup) {
	grupo := rg.Group("/blackjack")
	grupo.Use(middleware.JWTAuthMiddleware())

	ctrl := juegos.NewBlackjackController() 

	grupo.POST("/nueva", ctrl.CrearPartida)
	grupo.POST("/hit", ctrl.Handler(ctrl.Service.Hit))
	grupo.POST("/stand", ctrl.Handler(ctrl.Service.Stand))
	grupo.POST("/doblar", ctrl.Handler(ctrl.Service.Doblar))
	grupo.POST("/rendirse", ctrl.Handler(ctrl.Service.Rendirse))
	grupo.POST("/seguro", ctrl.Handler(ctrl.Service.Seguro))
	grupo.POST("/split", ctrl.Handler(ctrl.Service.Split))
	grupo.GET("/estado/:id_partida", ctrl.Estado)
}
