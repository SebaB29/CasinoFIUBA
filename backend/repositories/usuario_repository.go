package repositories

import (
	"casino/models"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db} // asignar correctamente
}

func (r *UsuarioRepository) Crear(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

// Busca un usuario por su ID
func (r *UsuarioRepository) ObtenerPorID(id uint) (*models.Usuario, error) {
	var u models.Usuario
	err := r.db.First(&u, id).Error
	return &u, err
}

// Devuelve la lista de todos los usuarios
func (r *UsuarioRepository) ObtenerTodos() ([]models.Usuario, error) {
	var lista []models.Usuario
	err := r.db.Find(&lista).Error
	return lista, err
}
