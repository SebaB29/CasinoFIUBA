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
type ResultadoIndividualDTO struct {
	TipoApuesta   string      `json:"tipo_apuesta"`
	MontoApostado float64     `json:"monto_apostado"`
	Ganancia      float64     `json:"ganancia"`
	Detalles      interface{} `json:"detalles"` // flexible según tipo
}

type RuletaResultadoUsuarioDTO struct {
	Mensaje       string                   `json:"message"`
	NumeroGanador int                      `json:"numero_ganador"`
	Color         string                   `json:"color"`
	Resultados    []ResultadoIndividualDTO `json:"resultados"`
}
