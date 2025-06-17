package main

import (
	"casino/config"
	"casino/db"
	"casino/routes"
	"math/rand"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Aleatoriedad global
	rand.Seed(time.Now().UnixNano()) 

	config.Load()
	// Establece la conexion con la base de datos
	db.ConectarDB()

	// Inicializa las rutas del servidor (con GIN)
	router := routes.SetupRoutes()

	// Obtiene el puerto desde variables de entorno o usar el default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicia el servidor en el puerto especificado (es el 8080 por defecto)
	router.Run(":" + port)
}
