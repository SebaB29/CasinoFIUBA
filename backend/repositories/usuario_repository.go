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
	Actualizar(usuario *models.Usuario) error
}

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (repository *UsuarioRepository) Crear(usuario *models.Usuario) error {
	return repository.db.Create(usuario).Error
}

func (repository *UsuarioRepository) ObtenerPorEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := repository.db.Where("email = ?", email).First(&usuario).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // no existe
	}
	return &usuario, err
}

func (repository *UsuarioRepository) ObtenerPorID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	err := repository.db.First(&usuario, id).Error
	return &usuario, err
}

func (repository *UsuarioRepository) ObtenerTodos() ([]models.Usuario, error) {
	var lista []models.Usuario
	err := repository.db.Find(&lista).Error
	return lista, err
}

func (repository *UsuarioRepository) Actualizar(usuario *models.Usuario) error {
	return repository.db.Save(usuario).Error
}
