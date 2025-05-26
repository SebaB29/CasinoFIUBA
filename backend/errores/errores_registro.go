package errores

import "errors"

var (
	ErrEmailYaExiste = errors.New("ya existe un usuario registrado con ese email")
	ErrMenorDeEdad   = errors.New("tienes que ser mayor de edad para registrarte")
	ErrGenerico      = errors.New("error al crear usuario")
)
