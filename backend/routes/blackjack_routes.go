package routes

import (
	ctrls "casino/controllers/juegos"
	dto "casino/dto/juegos"
	"casino/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func BlackjackRoutes(rg *gin.RouterGroup) {
	toGinH := func(fn func(uint, uint) (dto.BlackjackEstadoDTO, error)) func(uint, uint) (gin.H, error) {
		return func(idMesa uint, userID uint) (gin.H, error) {
			resp, err := fn(idMesa, userID)
			if err != nil {
				return nil, err
			}

			// Convertir el resultado a gin.H (Esto es necesario porque el controlador espera un gin.H para la respuesta)
			var m map[string]interface{}
			data, _ := json.Marshal(resp)
			_ = json.Unmarshal(data, &m)

			return gin.H(m), nil
		}
	}

	grupo := rg.Group("/blackjack")
	grupo.Use(middleware.JWTAuthMiddleware())

	ctrl := ctrls.NuevoBlackjackController()

	// Crea una nueva mesa de Blackjack
	grupo.POST("/mesa/nueva", ctrl.NuevaMesa)

	// Se une a una mesa existente (necesita id_mesa y apuesta inicial)
	grupo.POST("/mesa/unirse", ctrl.UnirseAMesa)

	// Obtiene el estado completo de la mesa
	grupo.GET("/mesa/estado/:id_mesa", ctrl.EstadoMesa)

	// Jugadas usando el wrapper
	grupo.POST("/mesa/hit", ctrl.Handler(toGinH(ctrl.Service.Hit)))
	grupo.POST("/mesa/stand", ctrl.Handler(toGinH(ctrl.Service.Stand)))
	grupo.POST("/mesa/doblar", ctrl.Handler(toGinH(ctrl.Service.Doblar)))
	grupo.POST("/mesa/rendirse", ctrl.Handler(toGinH(ctrl.Service.Rendirse)))
	grupo.POST("/mesa/seguro", ctrl.Handler(toGinH(ctrl.Service.Seguro)))
	grupo.POST("/mesa/split", ctrl.Handler(toGinH(ctrl.Service.Split)))

	// WebSocket
	grupo.GET("/mesa/ws", ctrl.JugarWS)
}
