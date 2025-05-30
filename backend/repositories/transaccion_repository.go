package repositories

import (
	"casino/models"

	"gorm.io/gorm"
)

type TransaccionRepositoryInterface interface {
	Crear(transaccion *models.Transaccion) error
}

type TransaccionRepository struct {
	db *gorm.DB
}

func NewTransaccionRepository(db *gorm.DB) *TransaccionRepository {
	return &TransaccionRepository{db: db}
}

func (r *TransaccionRepository) Crear(transaccion *models.Transaccion) error {
	return r.db.Create(transaccion).Error
}
