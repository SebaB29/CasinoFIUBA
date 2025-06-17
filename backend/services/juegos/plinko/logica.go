package plinko

import (
	dto "casino/dto/juegos"
	"math/rand"
)

const (
	niveles             = 8
	centro              = 4
	numPosiciones       = 9
	posicionMaximaIndex = numPosiciones - 1
	movIzquierda        = 0
	movDerecha          = 1
	DireccionIzquierda  = "I"
	DireccionDerecha    = "D"
)

var multiplicadores = [numPosiciones]float64{
	3.0, 1.4, 1.0, 0.8, 0.5, 0.8, 1.0, 1.4, 3.0,
}

// EjecutarPlinko simula una jugada de Plinko dado un monto y un generador aleatorio.
func EjecutarPlinko(monto float64, r *rand.Rand) dto.PlinkoResponseDTO {
	pos := centro
	// trayecto representa la secuencia de rebotes: "I" (izquierda), "D" (derecha)
	trayecto := make([]string, 0, niveles)

	for i := 0; i < niveles; i++ {
		direccion := r.Intn(2)
		if direccion == movIzquierda {
			trayecto = append(trayecto, DireccionIzquierda)
			if pos > 0 {
				pos--
			}
		} else {
			trayecto = append(trayecto, DireccionDerecha)
			if pos < posicionMaximaIndex {
				pos++
			}
		}
	}

	multiplicador := multiplicadores[pos]
	ganancia := monto * multiplicador

	return dto.PlinkoResponseDTO{
		PosicionFinal: pos,
		Multiplicador: multiplicador,
		Ganancia:      ganancia,
		Trayecto:      trayecto,
	}
}