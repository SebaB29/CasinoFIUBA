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
		UsuarioRoutes(v1)
		TransaccionRoutes(v1)
		BuscaminasRoutes(v1)
		VasosRoutes(v1)
		JuegosRoutes(v1)
	}

	return r
}
