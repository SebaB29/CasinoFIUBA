package dto

type TransaccionDTO struct {
	Monto float64 `json:"monto" binding:"required"`
	Tipo  string  `json:"tipo" binding:"required"`
}
