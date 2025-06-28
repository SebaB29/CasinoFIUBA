package models

import (
	"gorm.io/gorm"
)

type EstadoBlackjack string

const (
	EnCurso   EstadoBlackjack = "en_curso"
	Ganada    EstadoBlackjack = "ganada"
	Perdida   EstadoBlackjack = "perdida"
	Empatada  EstadoBlackjack = "empatada"
	Rendida   EstadoBlackjack = "rendida"
)

type MesaBlackjack struct {
	gorm.Model
	CartasBanca    string                  `json:"cartas_banca"`
	Mazo           string                  `json:"mazo"`
	Estado         string                  `json:"estado"`
	JugadorActual  int                     `gorm:"column:jugador_actual"`
	ManosJugadores []ManoJugadorBlackjack `gorm:"foreignKey:MesaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"manos_jugadores"`
}

type ManoJugadorBlackjack struct {
	gorm.Model
	MesaID       uint    `json:"mesa_id"`
	UserID       uint    `json:"user_id"`
	Apuesta      float64 `json:"apuesta"`
	ApuestaSplit float64 `gorm:"default:0" json:"apuesta_split"`
	Cartas       string  `json:"cartas"`
	CartasSplit  string  `json:"cartas_split"`
	ManoActual   int     `gorm:"default:1" json:"mano_actual"`
	Estado       string  `gorm:"default:'en_curso'" json:"estado"`
	Seguro       float64 `gorm:"default:0" json:"seguro"`
    ResultadoMano1 string `json:"resultado_mano_1"`
    ResultadoMano2 string `json:"resultado_mano_2"`
}
