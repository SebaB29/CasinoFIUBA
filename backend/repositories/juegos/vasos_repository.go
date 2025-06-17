package repositories

import (
	"casino/db"
	"casino/models"
)

func CrearPartidaVasos(partida *models.PartidaVasos) error {
	return db.DB.Create(partida).Error
}

func ObtenerPartidaVasosPorID(id uint) (*models.PartidaVasos, error) {
	var partida models.PartidaVasos
	err := db.DB.First(&partida, id).Error
	return &partida, err
}

func ActualizarPartidaVasos(partida *models.PartidaVasos) error {
	return db.DB.Save(partida).Error
}
