package dto

type RuletaRequestDTO struct {
	Monto       float64 `json:"monto" binding:"required"`
	TipoApuesta string  `json:"tipo_apuesta" binding:"required"`
	Numeros     []int   `json:"numeros,omitempty"`
	Docena      int     `json:"docena,omitempty"`
	Color       string  `json:"color,omitempty"`
	Paridad     string  `json:"paridad,omitempty"`
	AltoBajo    string  `json:"alto_bajo,omitempty"`
}

type RuletaResponseDTO struct {
	Mensaje       string  `json:"mensaje"`
	NumeroGanador int     `json:"numero_ganador"`
	ColorGanador  string  `json:"color"`
	MontoApostado float64 `json:"monto_apostado"`
	Ganancia      float64 `json:"ganancia"`
}
