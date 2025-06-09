package models

import "time"

type PartidaBuscaminas struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UsuarioID      *uint      `json:"usuario_id,omitempty"` 
	Usuario        Usuario    `json:"-" gorm:"foreignKey:UsuarioID"` 
	MontoApostado  float64    `json:"monto_apostado"`
	CantidadMinas  int        `json:"cantidad_minas"`
	CeldasAbiertas int        `json:"celdas_abiertas"`
	Estado         string     `json:"estado"` 
	PremioFinal    *float64   `json:"premio_final,omitempty"`
	FechaInicio    time.Time  `json:"fecha_inicio"`
	FechaFin       *time.Time `json:"fecha_fin,omitempty"`
}

func (PartidaBuscaminas) TableName() string {
	return "partidas_buscaminas"
}
