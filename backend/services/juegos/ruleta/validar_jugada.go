package ruleta

import (
	dto "casino/dto/juegos"
	"casino/errores"
	"fmt"
)

var validadores = map[string]func(dto.RuletaRequestDTO) error{
	"pleno":     validarPleno,
	"dividida":  validarDividida,
	"calle":     validarCalle,
	"cuadro":    validarCuadro,
	"docena":    validarDocena,
	"color":     validarColor,
	"paridad":   validarParidad,
	"alto_bajo": validarAltoBajo,
}

func ValidarJugada(jugada dto.RuletaRequestDTO) error {
	// Validaciones comunes previas
	for _, numero := range jugada.Numeros {
		if numero < NumeroMinimoRuleta || numero > NumeroMaximoRuleta {
			return fmt.Errorf("%w: %d", errores.ErrNumeroInvalido, numero)
		}
	}

	if fn, ok := validadores[jugada.TipoApuesta]; ok {
		return fn(jugada)
	}

	return fmt.Errorf("%w: tipo '%s'", errores.ErrApuestaInvalida, jugada.TipoApuesta)
}

// ----------- Validadores -----------

func validarPleno(jugada dto.RuletaRequestDTO) error {
	if len(jugada.Numeros) != CantidadNumerosPleno {
		return fmt.Errorf("%w: se requiere un único número para apuesta plena", errores.ErrFormatoInvalido)
	}
	return nil
}

func validarDividida(jugada dto.RuletaRequestDTO) error {
	if len(jugada.Numeros) != CantidadNumerosDividida || !sonAdyacentes(jugada.Numeros[0], jugada.Numeros[1]) {
		return fmt.Errorf("%w: los números deben ser adyacentes para una dividida", errores.ErrFormatoInvalido)
	}
	return nil
}

func validarCalle(jugada dto.RuletaRequestDTO) error {
	if len(jugada.Numeros) != CantidadNumerosCalle || !esCalleValida(jugada.Numeros) {
		return fmt.Errorf("%w: los números deben formar una calle válida", errores.ErrFormatoInvalido)
	}
	return nil
}

func validarCuadro(jugada dto.RuletaRequestDTO) error {
	if len(jugada.Numeros) != CantidadNumerosCuadro || !esCuadroValido(jugada.Numeros) {
		return fmt.Errorf("%w: los números deben formar un cuadro válido", errores.ErrFormatoInvalido)
	}
	return nil
}

func validarDocena(jugada dto.RuletaRequestDTO) error {
	if jugada.Docena < PrimeraDocena || jugada.Docena > TerceraDocena {
		return fmt.Errorf("%w: docena debe estar entre 1 y 3", errores.ErrFormatoInvalido)
	}
	return nil
}

func validarColor(jugada dto.RuletaRequestDTO) error {
	if !ColoresValidos[jugada.Color] {
		return fmt.Errorf("%w: debe ser 'rojo' o 'negro'", errores.ErrColorInvalido)
	}
	return nil
}

func validarParidad(jugada dto.RuletaRequestDTO) error {
	if !ParidadesValidas[jugada.Paridad] {
		return fmt.Errorf("%w: paridad debe ser 'par' o 'impar'", errores.ErrFormatoInvalido)
	}
	return nil
}

func validarAltoBajo(jugada dto.RuletaRequestDTO) error {
	if !AltoBajoValidos[jugada.AltoBajo] {
		return fmt.Errorf("%w: debe ser 'alto' o 'bajo'", errores.ErrFormatoInvalido)
	}
	return nil
}

// ----------- Funciones Auxiliares -----------

func abs(numero int) int {
	if numero < 0 {
		return -numero
	}
	return numero
}

func contieneCero(numeros []int) bool {
	for _, numero := range numeros {
		if numero == 0 {
			return true
		}
	}
	return false
}

func contieneTodosLosNumeros(numeros []int, buscados []int) bool {
	numerosEnLista := make(map[int]bool)
	for _, numero := range numeros {
		numerosEnLista[numero] = true
	}

	for _, numero := range buscados {
		if !numerosEnLista[numero] {
			return false
		}
	}
	return true
}

func obtenerPosicionEnTablero(numero int) (fila int, columna int, ok bool) {
	for i, filaTablero := range tableroRuleta {
		for j, numeroRuleta := range filaTablero {
			if numeroRuleta.Valor == numero {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}

func sonAdyacentes(numero1 int, numero2 int) bool {
	fila1, columna1, existe1 := obtenerPosicionEnTablero(numero1)
	fila2, columna2, existe2 := obtenerPosicionEnTablero(numero2)

	if !existe1 || !existe2 {
		return false
	}

	diferenciaFila := fila1 - fila2
	diferenciaColumna := columna1 - columna2
	return (diferenciaFila == 0 && abs(diferenciaColumna) == 1) || (abs(diferenciaFila) == 1 && diferenciaColumna == 0)
}

func esCalleValida(numeros []int) bool {
	if contieneCero(numeros) {
		return false
	}

	for _, fila := range tableroRuleta {
		valores := extraerValores(fila[:])
		if contieneTodosLosNumeros(valores, numeros) {
			return true
		}
	}

	return false
}

func esCuadroValido(numeros []int) bool {
	if contieneCero(numeros) {
		return false
	}

	for i := 0; i < filasEnTablero-1; i++ {
		for j := 0; j < columnasEnTablero-1; j++ {
			cuadro := []NumeroRuleta{
				tableroRuleta[i][j], tableroRuleta[i][j+1],
				tableroRuleta[i+1][j], tableroRuleta[i+1][j+1],
			}
			valores := extraerValores(cuadro)
			if contieneTodosLosNumeros(valores, numeros) {
				return true
			}
		}
	}

	return false
}

func extraerValores(numeros []NumeroRuleta) []int {
	valores := make([]int, len(numeros))
	for i, numeroRuleta := range numeros {
		valores[i] = numeroRuleta.Valor
	}

	return valores
}
