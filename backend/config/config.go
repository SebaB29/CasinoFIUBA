package config

import (
	"log"
	"os"
)

var JwtSecret []byte
var IsDevMode bool

func Load() {
	// Cargar clave JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET no está seteado en las variables de entorno")
	}
	JwtSecret = []byte(secret)

	// Cargar modo de entorno (dev o prod)
	env := os.Getenv("ENV")
	IsDevMode = (env == "dev")
	log.Printf("Modo actual: %s", env)
}
