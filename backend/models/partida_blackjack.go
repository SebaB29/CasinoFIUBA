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

type PartidaBlackjack struct {
    gorm.Model
    UserID          uint              `json:"user_id"`
    Apuesta         float64           `json:"apuesta"`
    ApuestaSplit    float64           `json:"apuesta_split"` 
    CartasJugador   string            `json:"cartas_jugador"`
    CartasJugadorSplit string         `json:"cartas_jugador_split"`
    CartasBanca     string            `json:"cartas_banca"`
    ManoActual      int               `json:"mano_actual"` 
    Estado          EstadoBlackjack   `json:"estado"`
    Mazo            string            `json:"mazo"`
    BlackjackNatural bool             `json:"blackjack_natural"`
    Doblar           bool             `json:"doblar"`
    Seguro           float64          `json:"seguro"`
}

