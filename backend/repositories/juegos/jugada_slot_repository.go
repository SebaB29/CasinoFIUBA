package repositories

import (
	"casino/models"

	"gorm.io/gorm"
)

type JugadaSlotRepositoryInterface interface {
	Crear(jugada *models.JugadaSlot) error
}

type JugadaSlotRepository struct {
	db *gorm.DB
}

func NewJugadaSlotRepository(db *gorm.DB) JugadaSlotRepositoryInterface {
	return &JugadaSlotRepository{db}
}

func (repo *JugadaSlotRepository) Crear(jugada *models.JugadaSlot) error {
	return repo.db.Create(jugada).Error
}
