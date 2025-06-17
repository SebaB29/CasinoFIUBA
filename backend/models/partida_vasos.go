package models

import "time"

type PartidaVasos struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `json:"user_id"`
	PosicionCorrecta  int       `json:"posicion_correcta"` // 0, 1 o 2
	PosicionElegida   *int      `json:"posicion_elegida"`  // nil si no hay jugada
	Estado            string    `json:"estado"`            // "en_curso", "ganada", "perdida"
	Apuesta           float64   `json:"apuesta"`
	FechaCreacion     time.Time `json:"fecha_creacion"`
}

func (PartidaVasos) TableName() string {
	return "partidas_vasos"
}
