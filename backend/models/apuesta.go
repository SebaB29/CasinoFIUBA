package models

import "time"

// Representa una transaccion de un usuario
type Apuesta struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Juego     string    `json:"juego"`
	Monto     int       `json:"monto"`
	Resultado string    `json:"resultado"`
	Fecha     time.Time `json:"fecha"`
}

func (Apuesta) TableName() string {
	return "apuestas"
}
