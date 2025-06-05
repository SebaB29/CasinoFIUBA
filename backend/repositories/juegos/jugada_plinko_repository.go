package repositories

import (
	"casino/models"

	"gorm.io/gorm"
)

type JugadaPlinkoRepositoryInterface interface {
	Crear(jugada *models.JugadaPlinko) error
}

type JugadaPlinkoRepository struct {
	db *gorm.DB
}

func NewJugadaPlinkoRepository(db *gorm.DB) *JugadaPlinkoRepository {
	return &JugadaPlinkoRepository{db: db}
}

func (r *JugadaPlinkoRepository) Crear(jugada *models.JugadaPlinko) error {
	return r.db.Create(jugada).Error
}
