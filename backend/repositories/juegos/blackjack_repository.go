package repositories

import (
    "casino/db"
    "casino/models"
    "gorm.io/gorm"
)

// CrearMesa crea una nueva mesa
func CrearMesa(mesa *models.MesaBlackjack) error {
    return db.DB.Create(mesa).Error
}

// ObtenerMesaConManos obtiene la mesa por ID y precarga las manos asociadas
func ObtenerMesaConManos(idMesa uint) (*models.MesaBlackjack, error) {
    var mesa models.MesaBlackjack
    err := db.DB.Preload("ManosJugadores").First(&mesa, idMesa).Error
    return &mesa, err
}

// ActualizarMesa actualiza la mesa completa (incluidas las manos)
func ActualizarMesa(mesa *models.MesaBlackjack) error {
    return db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(mesa).Error
}

// ActualizarManoJugador actualiza solo una mano específica
func ActualizarManoJugador(mano *models.ManoJugadorBlackjack) error {
    return db.DB.Save(mano).Error
}
