package juegos

type IniciarBlackjackDTO struct {
	Apuesta float64 `json:"apuesta" binding:"required"`
}

type JugadaBlackjackDTO struct {
	IDPartida uint   `json:"id_partida" binding:"required"`
	Accion    string `json:"accion" binding:"required,oneof=hit stand"`
}
