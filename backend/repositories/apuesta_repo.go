package repositories

import (
	"casino/models"
	"gorm.io/gorm"
)

type ApuestaRepository struct {
	db *gorm.DB
}

func NewApuestaRepository(db *gorm.DB) *ApuestaRepository {
	return &ApuestaRepository{db}
}

func (r *ApuestaRepository) Crear(apuesta *models.Apuesta) error {
	// TODO: guardar apuesta
	return r.db.Create(apuesta).Error
}

func (r *ApuestaRepository) ObtenerPorUsuario(userID uint) ([]models.Apuesta, error) {
	// TODO: obtener historial
	var apuestas []models.Apuesta
	err := r.db.Where("user_id = ?", userID).Find(&apuestas).Error
	return apuestas, err
}
