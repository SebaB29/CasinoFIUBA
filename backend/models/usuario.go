package models

import "time"

// Representa a un jugador del casino
type Usuario struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	Nombre          string        `json:"nombre"`
	Apellido        string        `json:"apellido"`
	FechaNacimiento time.Time     `json:"fecha_nacimiento"`
	Email           string        `json:"email"`
	Password        string        `json:"password"`
	Saldo           float64       `json:"saldo" gorm:"default:0"`
	Rol             string        `json:"rol" gorm:"default:user"`
	Transacciones   []Transaccion `json:"transacciones" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
