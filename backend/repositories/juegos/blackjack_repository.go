package repositories

import (
	"casino/db"
	"casino/models"
)

// CrearPartida guarda una nueva partida de blackjack en la base de datos
func CrearPartida(partida *models.PartidaBlackjack) error {
	return db.DB.Create(partida).Error
}

// ObtenerPartidaPorID busca una partida por ID
func ObtenerPartidaPorID(id uint) (*models.PartidaBlackjack, error) {
	var partida models.PartidaBlackjack
	err := db.DB.First(&partida, id).Error
	return &partida, err
}

// ActualizarPartida actualiza el estado de una partida existente
func ActualizarPartida(partida *models.PartidaBlackjack) error {
	return db.DB.Save(partida).Error
}
