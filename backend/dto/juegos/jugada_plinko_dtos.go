package dto

type PlinkoRequestDTO struct {
	Monto float64 `json:"monto" binding:"required"`
}

type PlinkoResponseDTO struct {
	PosicionFinal int      `json:"posicion_final"`
	Multiplicador float64  `json:"multiplicador"`
	Ganancia      float64  `json:"ganancia"`
	Trayecto      []string `json:"trayecto"`
}