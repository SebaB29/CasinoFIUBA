package repositories

import (
	"casino/models"
	"gorm.io/gorm"
)

type BuscaminasRepositoryInterface interface {
	Crear(partida *models.PartidaBuscaminas) error
	BuscarPorID(id uint) (*models.PartidaBuscaminas, error)
	Actualizar(partida *models.PartidaBuscaminas) error
}

type BuscaminasRepository struct {
	db *gorm.DB
}

func NewBuscaminasRepository(db *gorm.DB) *BuscaminasRepository {
	return &BuscaminasRepository{db: db}
}

func (r *BuscaminasRepository) Crear(partida *models.PartidaBuscaminas) error {
	return r.db.Create(partida).Error
}

func (r *BuscaminasRepository) BuscarPorID(id uint) (*models.PartidaBuscaminas, error) {
	var partida models.PartidaBuscaminas
	if err := r.db.First(&partida, id).Error; err != nil {
		return nil, err
	}
	return &partida, nil
}

func (r *BuscaminasRepository) Actualizar(partida *models.PartidaBuscaminas) error {
	return r.db.Save(partida).Error
}
