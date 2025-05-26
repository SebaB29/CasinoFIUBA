package services

import (
	"casino/db"
	"casino/dto"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	"errors"
)

type UsuarioService struct {
	repo *repositories.UsuarioRepository
}

func NewUsuarioService() *UsuarioService {
	repo := repositories.NewUsuarioRepository(db.DB) // usar db.DB global exportado
	return &UsuarioService{repo: repo}
}

func (s *UsuarioService) CrearUsuario(input dto.CrearUsuarioDTO) (*models.Usuario, error) {
	// Validación: edad mínima
	if input.Edad < 18 {
		return nil, errores.ErrMenorDeEdad
	}

	existe, _ := s.repo.ObtenerPorEmail(input.Email)
	if existe != nil {
		return nil, errores.ErrEmailYaExiste
	}

	usuario := models.Usuario{
		Nombre:   input.Nombre,
		Apellido: input.Apellido,
		Edad:     input.Edad,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := s.repo.Crear(&usuario); err != nil {
		return nil, errors.New("error al crear usuario")
	}

	return &usuario, nil
}
