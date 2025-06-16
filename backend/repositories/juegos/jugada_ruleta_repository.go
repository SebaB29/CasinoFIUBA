package repositories

import (
	"casino/models"

	"gorm.io/gorm"
)

type JugadaRuletaRepositoryInterface interface {
	Crear(jugada *models.JugadaRuleta) error
}

type JugadaRuletaRepository struct {
	db *gorm.DB
}

func NewJugadaRuletaRepository(db *gorm.DB) *JugadaRuletaRepository {
	return &JugadaRuletaRepository{db: db}
}

func (r *JugadaRuletaRepository) Crear(jugada *models.JugadaRuleta) error {
	return r.db.Create(jugada).Error
}