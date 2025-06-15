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
