package services

import (
	dto "casino/dto/juegos"
	"casino/errores"
	"fmt"
)

func ValidarApuesta(jugada dto.RuletaRequestDTO) error {
	for _, n := range jugada.Numeros {
		if n < 0 || n > 36 {
			return fmt.Errorf("%w: %d", errores.ErrNumeroInvalido, n)
		}
	}

	switch jugada.TipoApuesta {
	case "pleno":
		if len(jugada.Numeros) != 1 {
			return fmt.Errorf("%w: se requiere un único número para apuesta plena", errores.ErrFormatoInvalido)
		}
	case "dividida":
		if len(jugada.Numeros) != 2 || !SonAdyacentes(jugada.Numeros[0], jugada.Numeros[1]) {
			return fmt.Errorf("%w: los números deben ser adyacentes para una dividida", errores.ErrFormatoInvalido)
		}
	case "calle":
		if len(jugada.Numeros) != 3 || !EsCalleValida(jugada.Numeros) {
			return fmt.Errorf("%w: los números deben formar una calle válida", errores.ErrFormatoInvalido)
		}
	case "cuadro":
		if len(jugada.Numeros) != 4 || !EsCuadroValido(jugada.Numeros) {
			return fmt.Errorf("%w: los números deben formar un cuadro válido", errores.ErrFormatoInvalido)
		}
	case "docena":
		if jugada.Docena < 1 || jugada.Docena > 3 {
			return fmt.Errorf("%w: docena debe estar entre 1 y 3", errores.ErrFormatoInvalido)
		}
	case "color":
		if jugada.Color != "rojo" && jugada.Color != "negro" {
			return fmt.Errorf("%w: debe ser 'rojo' o 'negro'", errores.ErrColorInvalido)
		}
	case "paridad":
		if jugada.Paridad != "par" && jugada.Paridad != "impar" {
			return fmt.Errorf("%w: paridad debe ser 'par' o 'impar'", errores.ErrFormatoInvalido)
		}
	case "alto_bajo":
		if jugada.AltoBajo != "alto" && jugada.AltoBajo != "bajo" {
			return fmt.Errorf("%w: debe ser 'alto' o 'bajo'", errores.ErrFormatoInvalido)
		}
	default:
		return fmt.Errorf("%w: tipo '%s'", errores.ErrApuestaInvalida, jugada.TipoApuesta)
	}
	return nil
}
