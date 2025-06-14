package routes

import (
	"casino/controllers/juegos"
	"github.com/gin-gonic/gin"
)

func RegistrarRutasBlackjack(r *gin.Engine) {
	grupo := r.Group("/blackjack")
	{
		grupo.POST("/nueva", juegos.CrearPartidaBlackjack)
		grupo.POST("/hit", juegos.HitBlackjack)
		grupo.POST("/stand", juegos.StandBlackjack)
		grupo.GET("/estado/:id_partida", juegos.ObtenerEstadoBlackjack)
	}
}
