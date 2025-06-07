package routes

import (
	"casino/controllers"
	"github.com/gin-gonic/gin"
)

func BuscaminasRoutes(rg *gin.RouterGroup) {
	ruta := rg.Group("/buscaminas")
	ruta.POST("/nueva", controllers.CrearPartidaBuscaminas)
	ruta.POST("/abrir", controllers.AbrirCeldaBuscaminas)
	ruta.POST("/retirarse", controllers.RetirarseBuscaminas)

}
