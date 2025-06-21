package dto

type IniciarBlackjackDTO struct {
	Apuesta float64 `json:"apuesta" binding:"required"`
}

type JugadaBlackjackDTO struct {
	IDPartida uint   `json:"id_partida" binding:"required"`
}

type SplitBlackjackDTO struct {
    IDPartida uint `json:"id_partida" binding:"required"`
}

type SeguroBlackjackDTO struct {
    IDPartida uint `json:"id_partida" binding:"required"`
}