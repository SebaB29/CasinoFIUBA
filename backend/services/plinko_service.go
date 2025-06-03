package services

import (
	"math/rand"
)

type PlinkoResultado struct {
	PosicionFinal int     `json:"posicion_final"`
	Multiplicador float64 `json:"multiplicador"`
	Ganancia      float64 `json:"ganancia"`
}

// Esta función representa una jugada de prueba con monto fijo
func JugarPlinko(monto float64) PlinkoResultado {
	const niveles = 8

	// Simulación de rebotes
	pos := 0
	for i := 0; i < niveles; i++ {
		direccion := rand.Intn(2) // 0 = izquierda, 1 = derecha
		pos += direccion
	}

	// Multiplicadores típicos de un tablero de 9 posiciones
	multiplicadores := []float64{0.5, 0.8, 1.0, 1.4, 3.0, 1.4, 1.0, 0.8, 0.5}

	multiplicador := multiplicadores[pos]
	ganancia := float64(monto) * multiplicador

	return PlinkoResultado{
		PosicionFinal: pos,
		Multiplicador: multiplicador,
		Ganancia:      ganancia,
	}
}
