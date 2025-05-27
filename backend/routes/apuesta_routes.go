package routes

import (
	"casino/controllers"

	"github.com/gin-gonic/gin"
)

func RegistroApuestaRoutes(rg *gin.RouterGroup) {
	apuestas := rg.Group("/apuestas")

	apuestas.POST("/", controllers.CrearApuesta)
}
