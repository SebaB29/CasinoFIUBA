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
)

type PartidaBlackjack struct {
    gorm.Model
    UserID        uint              `json:"user_id"`
    Apuesta       float64           `json:"apuesta"`
    CartasJugador string            `json:"cartas_jugador"`
    CartasBanca   string            `json:"cartas_banca"`  
    Estado        EstadoBlackjack   `json:"estado"`
    Mazo          string            `json:"mazo"` 
}
