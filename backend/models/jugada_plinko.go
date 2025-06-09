package models

import "time"

type JugadaPlinko struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UsuarioID     uint      `gorm:"not null" json:"usuario_id"`
	Usuario       Usuario   `json:"-" gorm:"foreignKey:UsuarioID"` // Relación con Usuario, opcional en respuesta JSON
	MontoApostado float64   `json:"monto_apostado"`
	Multiplicador float64   `json:"multiplicador"`
	Ganancia      float64   `json:"ganancia"`
	PosicionFinal int       `json:"posicion_final"`
	Fecha         time.Time `json:"fecha"`
}

func (JugadaPlinko) TableName() string {
	return "jugadas_plinko"
}
