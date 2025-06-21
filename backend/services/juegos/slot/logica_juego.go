package slot

import (
	"math/rand"
	"time"
)

var simbolos = []string{"🍒", "🍋", "🔔", "⭐", "🍇"}

func generarResultadoSlot(apuesta float64) ([][]string, float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	rondas := [][]string{}
	totalGanancia := 0.0

	for i := 0; i < 3; i++ {
		fila := []string{}
		contador := map[string]int{}
		for j := 0; j < 3; j++ {
			simbolo := simbolos[r.Intn(len(simbolos))]
			fila = append(fila, simbolo)
			contador[simbolo]++
		}
		rondas = append(rondas, fila)
		totalGanancia += calcularGanancia(apuesta, contador)
	}

	return rondas, totalGanancia
}

func calcularGanancia(apuesta float64, contador map[string]int) float64 {
	ganancia := 0.0
	for _, count := range contador {
		if count == 3 {
			ganancia += apuesta * 5
		}
	}
	return ganancia
}

func generarMensajeSlot(ganancia float64) string {
	if ganancia > 0 {
		return "¡Ganaste!"
	}
	return "Suerte la próxima..."
}
