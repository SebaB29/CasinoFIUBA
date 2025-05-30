package errores

import "errors"

var (
	ErrEmailYaExiste = errors.New("ya existe un usuario registrado con ese correo electrónico")
	ErrMenorDeEdad   = errors.New("debes ser mayor de edad para registrarte")
	ErrGenerico      = errors.New("no se pudo crear el usuario, intenta nuevamente más tarde")
)
