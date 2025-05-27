package models

// Representa a un jugador del casino
type Usuario struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Edad     uint   `json:"edad"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
