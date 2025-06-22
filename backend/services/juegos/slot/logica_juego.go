package slot

import (
	"math/rand"
	"time"
)

type Simbolo struct {
	Emoji         string
	Multiplicador float64
	Probabilidad  int
}

var simbolos = []Simbolo{
	{"🍒", 2.0, 40},
	{"🍋", 3.0, 25},
	{"🔔", 5.0, 20},
	{"⭐", 8.0, 10},
	{"💎", 15.0, 5},
}

func generarResultadoSlot(apuesta float64) ([][]string, float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	rondas := [][]string{}
	totalGanancia := 0.0

	for i := 0; i < 3; i++ {
		fila := []string{}
		contador := map[string]int{}

		for j := 0; j < 3; j++ {
			simbolo := obtenerSimboloPonderado(r)
			fila = append(fila, simbolo.Emoji)
			contador[simbolo.Emoji]++
		}

		totalGanancia += calcularGanancia(apuesta, contador)
		rondas = append(rondas, fila)
	}

	return rondas, totalGanancia
}

func obtenerSimboloPonderado(r *rand.Rand) Simbolo {
	total := 0
	for _, simbolo := range simbolos {
		total += simbolo.Probabilidad
	}

	valor := r.Intn(total)
	acumulado := 0

	for _, simbolo := range simbolos {
		acumulado += simbolo.Probabilidad
		if valor < acumulado {
			return simbolo
		}
	}

	return simbolos[0]
}

func calcularGanancia(apuesta float64, contador map[string]int) float64 {
	for emoji, count := range contador {
		if count == 3 {
			multiplicador := obtenerMultiplicadorPorEmoji(emoji)
			return apuesta * multiplicador
		}
	}
	return 0
}

func obtenerMultiplicadorPorEmoji(emoji string) float64 {
	for _, simbolo := range simbolos {
		if simbolo.Emoji == emoji {
			return simbolo.Multiplicador
		}
	}
	return 0
}

func generarMensajeSlot(ganancia float64) string {
	if ganancia > 0 {
		return "¡Ganaste!"
	}
	return "Suerte la próxima..."
}
