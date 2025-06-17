package dto

type CrearPartidaVasosDTO struct {
	Apuesta float64 `json:"apuesta" binding:"required,gt=0"`
}

type RealizarJugadaVasosDTO struct {
	IDPartida uint `json:"id_partida" binding:"required"`
	Eleccion  *int  `json:"eleccion" binding:"required,min=0,max=2"` // Solo 0, 1 o 2
}
