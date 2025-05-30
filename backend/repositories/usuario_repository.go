package repositories

import (
	"casino/models"

	"gorm.io/gorm"
)

type UsuarioRepositoryInterface interface {
	Crear(usuario *models.Usuario) error
	ObtenerPorEmail(email string) (*models.Usuario, error)
	ObtenerPorID(id uint) (*models.Usuario, error)
	ObtenerTodos() ([]models.Usuario, error)
}

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) Crear(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *UsuarioRepository) ObtenerPorEmail(email string) (*models.Usuario, error) {
	var u models.Usuario
	err := r.db.Where("email = ?", email).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, err
}

func (r *UsuarioRepository) ObtenerPorID(id uint) (*models.Usuario, error) {
	var u models.Usuario
	err := r.db.First(&u, id).Error
	return &u, err
}

func (r *UsuarioRepository) ObtenerTodos() ([]models.Usuario, error) {
	var lista []models.Usuario
	err := r.db.Find(&lista).Error
	return lista, err
}
