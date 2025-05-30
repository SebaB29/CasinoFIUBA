package errores

import "errors"

var (
	ErrUsuarioNoEncontrado = errors.New("el usuario no está registrado")
	ErrPasswordIncorrecta  = errors.New("la contraseña ingresada es incorrecta")
)
