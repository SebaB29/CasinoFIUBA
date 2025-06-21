package blackjack

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
    CantidadMazos = 6
)

// Mazo completo de 52 cartas (4 palos)
var mazoCompleto = []string{
	"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
	"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
	"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
	"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
}

func MezclarMazo(numMazos int) []string {
    baraja := []string{}
    for i := 0; i < numMazos; i++ {
        baraja = append(baraja, mazoCompleto...)
    }
    rand.Shuffle(len(baraja), func(i, j int) {
        baraja[i], baraja[j] = baraja[j], baraja[i]
    })
    return baraja
}

func TomarCarta(mazo []string) (string, []string) {
	if len(mazo) == 0 {
		return "", mazo
	}
	return mazo[0], mazo[1:]
}

func CalcularValor(cartas []string) int {
	total := 0
	ases := 0

	for _, carta := range cartas {
		switch carta {
		case "J", "Q", "K":
			total += 10
		case "A":
			total += 11
			ases++
		default:
			var valor int
			fmt.Sscanf(carta, "%d", &valor)
			total += valor
		}
	}

	for total > 21 && ases > 0 {
		total -= 10
		ases--
	}

	return total
}

func CartasToString(cartas []string) string {
	return strings.Join(cartas, "-")
}

func StringToCartas(cadena string) []string {
	if cadena == "" {
		return []string{}
	}
	return strings.Split(cadena, "-")
}

func TieneBlackjack(cartas []string) bool {
	if len(cartas) != 2 {
		return false
	}

	normalizar := func(c string) string {
		return strings.ToUpper(strings.TrimSpace(c))
	}

	c1 := normalizar(cartas[0])
	c2 := normalizar(cartas[1])

	esAs := func(c string) bool { return c == "A" }
	esDiez := func(c string) bool { return c == "10" || c == "J" || c == "Q" || c == "K" }

	return (esAs(c1) && esDiez(c2)) || (esAs(c2) && esDiez(c1))
}

func Es17Blando(cartas []string) bool {
	total := 0
	ases := 0

	for _, carta := range cartas {
		switch carta {
		case "A":
			total += 11
			ases++
		case "J", "Q", "K":
			total += 10
		default:
			var valor int
			fmt.Sscanf(carta, "%d", &valor)
			total += valor
		}
	}

	return total == 17 && ases > 0
}

// EvaluarResultado determina el resultado de la partida de Blackjack
func EvaluarResultado(valorJugador, valorBanca int) string {
	if valorJugador > 21 {
		return "perdida"
	}
	if valorBanca > 21 {
		return "ganada"
	}
	if valorJugador > valorBanca {
		return "ganada"
	}
	if valorJugador < valorBanca {
		return "perdida"
	}
	return "empatada"
}
