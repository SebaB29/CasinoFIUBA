package dto

type CrearUsuarioDTO struct {
	Nombre          string `json:"nombre" binding:"required"`
	Apellido        string `json:"apellido"`
	FechaNacimiento string `json:"fecha_nacimiento" binding:"required"` // Formato: YYYY-MM-DD
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
}

