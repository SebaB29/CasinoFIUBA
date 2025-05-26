package errores

import "errors"

var (
	ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")
	ErrPasswordIncorrecta  = errors.New("contraseña incorrecta")
)
