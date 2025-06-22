package models

import "time"

type JugadaSlot struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UsuarioID     uint      `json:"usuario_id" gorm:"not null"`
	Usuario       Usuario   `json:"-" gorm:"foreignKey:UsuarioID"`
	MontoApostado float64   `json:"monto_apostado"`
	Ganancia      float64   `json:"ganancia"`
	Rondas        string    `json:"rondas"`
	Fecha         time.Time `json:"fecha"`
}

func (JugadaSlot) TableName() string {
	return "jugadas_slot"
}
