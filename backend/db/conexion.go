package db

import (
	"fmt"
	"log"
	"os"

	"casino/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConectarDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Faltan variables de entorno necesarias para la conexión a la DB")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	DB = db
	log.Println("Base de datos conectada correctamente")

	// Se eliminan las tablas para poder testear bien las cosas. LUEGO SACAR
	db.Migrator().DropTable(&models.Usuario{})
	db.Migrator().DropTable(&models.Transaccion{})
	db.Migrator().DropTable(&models.JugadaPlinko{})

	// Se crean las tablas de la BD
	// Migrar modelo Usuario (crea tabla si no existe)
	if err := DB.AutoMigrate(&models.Usuario{}); err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}
	if err := DB.AutoMigrate(&models.Transaccion{}); err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}
	if err := DB.AutoMigrate(&models.JugadaPlinko{}); err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}
}
