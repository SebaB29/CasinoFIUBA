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
	CrearUsuario(input dto.CrearUsuarioDTO) (*models.Usuario, error)
	Login(input dto.LoginDTO) (*models.Usuario, string, error)
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

func (service *UsuarioService) CrearUsuario(input dto.CrearUsuarioDTO) (*models.Usuario, error) {
	fechaNacimiento, err := time.Parse(FormatoFechaNacimiento, input.FechaNacimiento)
	if err != nil {
		return nil, errores.ErrFormatoFechaInvalido
	}

	edad := calcularEdad(fechaNacimiento)
	if edad < EdadMinimaPermitida {
		return nil, errores.ErrMenorDeEdad
	}

	existe, _ := service.repository.ObtenerPorEmail(input.Email)
	if existe != nil {
		return nil, errores.ErrEmailYaExiste
	}

	usuario := models.Usuario{
		Nombre:          input.Nombre,
		Apellido:        input.Apellido,
		FechaNacimiento: fechaNacimiento,
		Email:           input.Email,
		Password:        input.Password,
	}

	if err := service.repository.Crear(&usuario); err != nil {
		return nil, errores.ErrGenerico
	}

	return &usuario, nil
}

func calcularEdad(nacimiento time.Time) int {
	hoy := time.Now()
	edad := hoy.Year() - nacimiento.Year()
	if hoy.YearDay() < nacimiento.YearDay() {
		edad--
	}
	return edad
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

// 	NUEVA FUNCIONALIDAD ObtenerPorID busca un usuario por su ID y devuelve un error si no se encuentra
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
