package services

import (
	"casino/db"
	"casino/dto"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	"casino/utils"
)

const EdadMinimaPermitida = 18

type UsuarioServiceInterface interface {
	CrearUsuario(input dto.CrearUsuarioDTO) (*models.Usuario, error)
	Login(input dto.LoginDTO) (*models.Usuario, string, error)
}

type UsuarioService struct {
	repository repositories.UsuarioRepositoryInterface
}

func NewUsuarioService() *UsuarioService {
	repository := repositories.NewUsuarioRepository(db.DB)
	return &UsuarioService{repository: repository}
}

func (service *UsuarioService) CrearUsuario(input dto.CrearUsuarioDTO) (*models.Usuario, error) {
	if input.Edad < EdadMinimaPermitida {
		return nil, errores.ErrMenorDeEdad
	}

	existe, _ := service.repository.ObtenerPorEmail(input.Email)
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

	if err := service.repository.Crear(&usuario); err != nil {
		return nil, errores.ErrGenerico
	}

	return &usuario, nil
}

func (service *UsuarioService) Login(input dto.LoginDTO) (*models.Usuario, string, error) {
	usuario, err := service.repository.ObtenerPorEmail(input.Email)
	if err != nil || usuario == nil {
		return nil, "", errores.ErrUsuarioNoEncontrado
	}

	if usuario.Password != input.Password {
		return nil, "", errores.ErrPasswordIncorrecta
	}

	token, err := utils.GenerateToken(usuario.ID)
	if err != nil {
		return nil, "", err
	}

	return usuario, token, nil
}
