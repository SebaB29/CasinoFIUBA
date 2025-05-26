package routes

import (
	"casino/controllers"
	"github.com/gin-gonic/gin"
)

//Inicializa el router y define todas las rutas de la api
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Ruta base para verificar que el servidor esta activo
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"mensaje": "🎰 Casino API funcionando"})
	})

	// Rutas de usuarios
	r.POST("/usuarios", controllers.CrearUsuario)
	r.GET("/usuarios", controllers.ObtenerUsuarios)

	// Rutas de apuestas (provisoria)
	r.POST("/apuestas", controllers.CrearApuesta)

	return r
}
