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
	MontoApostado float64 `json:"monto_apostado"`
	TipoApuesta   string  `json:"tipo_apuesta"`
	Numeros       []int   `json:"numeros,omitempty"`
	Docena        int     `json:"docena,omitempty"`
	Color         string  `json:"color,omitempty"`
	Paridad       string  `json:"paridad,omitempty"`
	AltoBajo      string  `json:"alto_bajo,omitempty"`
	NumeroGanador int     `json:"numero_ganador"`
	ColorGanador  string  `json:"color_ganador"`
	Multiplicador float64 `json:"multiplicador"`
	Ganancia      float64 `json:"ganancia"`
}
