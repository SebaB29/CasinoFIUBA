package services

import (
	"casino/db"
	"casino/dto"
	"casino/errores"
	"casino/models"
	"casino/repositories"
	"casino/utils"
	"time"
)

const EdadMinimaPermitida = 18
const FormatoFechaNacimiento = "2006-01-02"

type UsuarioServiceInterface interface {
	CrearUsuario(request dto.RegistroUsuarioRequestDTO) (*dto.RegistroUsuarioResponseDTO, error)
	Login(request dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
	ObtenerPorID(id uint) (*models.Usuario, error)
	ObtenerTodos() ([]models.Usuario, error)
}

type UsuarioService struct {
	repository repositories.UsuarioRepositoryInterface
}

func NewUsuarioService() *UsuarioService {
	repository := repositories.NewUsuarioRepository(db.DB)
	return &UsuarioService{repository: repository}
}

func (service *UsuarioService) CrearUsuario(request dto.RegistroUsuarioRequestDTO) (*dto.RegistroUsuarioResponseDTO, error) {
	fechaNacimiento, err := time.Parse("2006-01-02", request.FechaNacimiento)
	if err != nil {
		return nil, errores.ErrFormatoFechaInvalido
	}

	if esMenorDeEdad(fechaNacimiento) {
		return nil, errores.ErrMenorDeEdad
	}

	existe, _ := service.repository.ObtenerPorEmail(request.Email)
	if existe != nil {
		return nil, errores.ErrEmailYaExiste
	}

	usuario := models.Usuario{
		Nombre:          request.Nombre,
		Apellido:        request.Apellido,
		FechaNacimiento: fechaNacimiento,
		Email:           request.Email,
		Password:        request.Password,
	}

	if err := service.repository.Crear(&usuario); err != nil {
		return nil, errores.ErrGenerico
	}

	response := dto.RegistroUsuarioResponseDTO{
		ID:              usuario.ID,
		Nombre:          usuario.Nombre,
		Apellido:        usuario.Apellido,
		FechaNacimiento: fechaNacimiento,
		Email:           usuario.Email,
		Saldo:           usuario.Saldo,
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

	token, err := utils.GenerateToken(usuario.ID, usuario.Rol)
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

// NUEVA FUNCIONALIDAD ObtenerPorID busca un usuario por su ID y devuelve un error si no se encuentra
func (service *UsuarioService) ObtenerPorID(id uint) (*models.Usuario, error) {
	usuario, err := service.repository.ObtenerPorID(id)
	if err != nil {
		return nil, errores.ErrUsuarioNoEncontrado
	}
	return usuario, nil
}

// NUEVA FUNCIONALIDAD ObtenerTodos devuelve todos los usuarios registrados en el sistema
func (service *UsuarioService) ObtenerTodos() ([]models.Usuario, error) {
	return service.repository.ObtenerTodos()
}

func esMenorDeEdad(nacimiento time.Time) bool {
	hoy := time.Now()
	edad := hoy.Year() - nacimiento.Year()
	if hoy.Month() < nacimiento.Month() || (hoy.Month() == nacimiento.Month() && hoy.Day() < nacimiento.Day()) {
		edad--
	}
	return edad < EdadMinimaPermitida
}
