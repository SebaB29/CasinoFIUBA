package config

import (
	"log"
	"os"
)

var JwtSecret []byte

func Load() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET no está seteado en las variables de entorno")
	}
	JwtSecret = []byte(secret)
}