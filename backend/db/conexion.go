package db

import (
	"casino/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Establece la conexión a PostgreSQL y migra los modelos
func ConectarDB() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN no definido")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	if err := db.AutoMigrate(&models.Usuario{}, &models.Apuesta{}); err != nil {
		log.Fatalf("Error al migrar modelos: %v", err)
	}

	DB = db
}
