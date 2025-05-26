package models

// Representa a un jugador del casino
type Usuario struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
	Saldo  int    `json:"saldo"`
}
