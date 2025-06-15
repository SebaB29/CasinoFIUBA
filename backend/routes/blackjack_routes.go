package routes

import (
	juegos "casino/controllers/juegos"
	"github.com/gin-gonic/gin"
	"casino/middleware"
)

func BlackjackRoutes(rg *gin.RouterGroup) {
	grupo := rg.Group("/blackjack")
	grupo.Use(middleware.JWTAuthMiddleware())

	grupo.POST("/nueva", juegos.CrearPartidaBlackjack)
	grupo.POST("/hit", juegos.HitBlackjack)
	grupo.POST("/stand", juegos.StandBlackjack)
	grupo.GET("/estado/:id_partida", juegos.ObtenerEstadoBlackjack)
}
