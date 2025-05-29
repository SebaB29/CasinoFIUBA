package dto

type TransaccionDTO struct {
	Monto float64 `json:"monto" binding:"required,gt=0"`
	Tipo  string  `json:"tipo" binding:"required"`
}
