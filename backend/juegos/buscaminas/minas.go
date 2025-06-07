package buscaminas

import (
	"math/rand"
	"time"
)

// GenerarMinas coloca minas aleatorias en el tablero, sin duplicados
func (t *Tablero) GenerarMinas(cantidad int) {
	rand.Seed(time.Now().UnixNano())

	totalCeldas := t.Filas * t.Columnas
	if cantidad > totalCeldas {
		cantidad = totalCeldas
	}

	posiciones := make(map[int]bool)

	for len(posiciones) < cantidad {
		index := rand.Intn(totalCeldas)
		if !posiciones[index] {
			posiciones[index] = true
			fila := index / t.Columnas
			col := index % t.Columnas
			t.Celdas[fila][col].TieneMina = true
		}
	}
}
