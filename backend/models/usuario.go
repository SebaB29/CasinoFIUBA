package models

import "time"

// Representa a un jugador del casino
type Usuario struct {
<<<<<<< HEAD
	ID       uint   `json:"id" gorm:"primaryKey"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Edad     uint   `json:"edad"`
	Email    string `json:"email"`
	Password string `json:"password"`
=======
	ID              uint      `json:"id" gorm:"primaryKey"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
>>>>>>> 4948393 (feat: mejora en la validacion de edad por fecha de nacimiento. Comienzo de desarrollo para un get all users y users por id)
}
