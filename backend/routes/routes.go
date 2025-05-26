package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes inicializa el router y todas las rutas principales
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Ruta de prueba o estado del servidor
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensaje": "🎰 Casino API funcionando"})
	})

	// Grupo de rutas versión 1
	v1 := r.Group("")
	{
		RegistroUsuarioRoutes(v1)
		//RegistroApuestaRoutes(v1)
	}

	return r
}
