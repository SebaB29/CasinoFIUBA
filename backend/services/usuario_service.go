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
	CrearUsuario(request dto.RegistroUsuarioRequestDTO) (*dto.RegistroUsuarioResponseDTO, error)
	Login(request dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
}

type UsuarioService struct {
	repository repositories.UsuarioRepositoryInterface
}

func NewUsuarioService() *UsuarioService {
	repository := repositories.NewUsuarioRepository(db.DB)
	return &UsuarioService{repository: repository}
}

func (service *UsuarioService) CrearUsuario(request dto.RegistroUsuarioRequestDTO) (*dto.RegistroUsuarioResponseDTO, error) {
	if request.Edad < EdadMinimaPermitida {
		return nil, errores.ErrMenorDeEdad
	}

	existe, _ := service.repository.ObtenerPorEmail(request.Email)
	if existe != nil {
		return nil, errores.ErrEmailYaExiste
	}

	usuario := models.Usuario{
		Nombre:   request.Nombre,
		Apellido: request.Apellido,
		Edad:     request.Edad,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := service.repository.Crear(&usuario); err != nil {
		return nil, errores.ErrGenerico
	}

	response := dto.RegistroUsuarioResponseDTO{
		ID:       usuario.ID,
		Nombre:   usuario.Nombre,
		Apellido: usuario.Apellido,
		Edad:     usuario.Edad,
		Email:    usuario.Email,
		Saldo:    usuario.Saldo,
	}

	return &response, nil
}

func (service *UsuarioService) Login(request dto.LoginRequestDTO) (*dto.LoginResponseDTO, error) {
	usuario, err := service.repository.ObtenerPorEmail(request.Email)
	if err != nil || usuario == nil {
		return nil, errores.ErrUsuarioNoEncontrado
	}

	if usuario.Password != request.Password {
		return nil, errores.ErrPasswordIncorrecta
	}

	token, err := utils.GenerateToken(usuario.ID)
	if err != nil {
		return nil, err
	}

	response := dto.LoginResponseDTO{
		ID:     usuario.ID,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
		Token:  token,
	}

	return &response, nil
}
