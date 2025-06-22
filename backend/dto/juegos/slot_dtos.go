// dto/juegos/slot_dto.go

package dto

type SlotRequestDTO struct {
	Monto float64 `json:"monto"`
}

type SlotResponseDTO struct {
	Rondas   [][]string `json:"rondas"`
	Ganancia float64    `json:"ganancia"`
	Mensaje  string     `json:"mensaje"`
}
