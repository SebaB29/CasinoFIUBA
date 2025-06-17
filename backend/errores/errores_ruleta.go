package errores

import "errors"

var (
	ErrColorInvalido   = errors.New("color inválido")
	ErrNumeroInvalido  = errors.New("número fuera del rango 0-36")
	ErrFormatoInvalido = errors.New("formato de apuesta inválido")
	ErrApuestaInvalida = errors.New("tipo de apuesta no válida o datos incorrectos")
)
