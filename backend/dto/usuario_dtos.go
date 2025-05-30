package dto

type CrearUsuarioDTO struct {
	Nombre           string `json:"nombre" binding:"required"`
	Apellido         string `json:"apellido"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required"`
	FechaNacimiento  string `json:"fecha_nacimiento" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
