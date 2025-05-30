// dto/usuario_dtos.go
package dto

import "time"

// ---------------- REQUEST ----------------

type RegistroUsuarioRequestDTO struct {
	Nombre          string `json:"nombre" binding:"required"`
	Apellido        string `json:"apellido"`
	FechaNacimiento string `json:"fecha_nacimiento" binding:"required" time_format:"2006-01-02"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ---------------- RESPONSE ----------------

type RegistroUsuarioResponseDTO struct {
	ID              uint      `json:"id"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	FechaNacimiento time.Time `json:"fecha_nacimiento" time_format:"2006-01-02"`
	Email           string    `json:"email"`
	Saldo           float64   `json:"saldo"`
}

type LoginResponseDTO struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}
