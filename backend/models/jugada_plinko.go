package models

import "time"

type JugadaPlinko struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UsuarioID     uint      `json:"usuario_id" gorm:"not null"`
	Usuario       Usuario   `json:"-" gorm:"foreignKey:UsuarioID"`
	MontoApostado float64   `json:"monto_apostado"`
	Multiplicador float64   `json:"multiplicador"`
	Ganancia      float64   `json:"ganancia"`
	PosicionFinal int       `json:"posicion_final"`
	Trayecto      string    `json:"trayecto" gorm:"type:text"`
	Fecha         time.Time `json:"fecha"`
}

func (JugadaPlinko) TableName() string {
	return "jugadas_plinko"
}