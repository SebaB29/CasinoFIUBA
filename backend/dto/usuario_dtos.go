// dto/usuario_dtos.go
package dto

type CrearUsuarioDTO struct {
	Nombre   string `json:"nombre" binding:"required"`
	Apellido string `json:"apellido"`
	Edad     uint   `json:"edad" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
