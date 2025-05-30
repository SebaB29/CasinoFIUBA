package models

import "time"

type Transaccion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UsuarioID uint      `gorm:"not null" json:"user_id"`
	Usuario   Usuario   `json:"-" gorm:"foreignKey:UsuarioID"`
	Monto     float64   `gorm:"not null" json:"monto"`
	Tipo      string    `gorm:"type:varchar(20);not null" json:"tipo"`
	CreatedAt time.Time `json:"fecha"`
}

// TableName indica a GORM que esta estructura usa la tabla "transacciones"
func (Transaccion) TableName() string {
	return "transacciones"
}
