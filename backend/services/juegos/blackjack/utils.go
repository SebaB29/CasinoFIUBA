package blackjack

import (
	"fmt"
	"math/rand"
	"strings"
	"time" // Importar time para inicializar rand
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

func init() {
	// Inicializa la semilla del generador de números aleatorios una vez
	rand.Seed(time.Now().UnixNano())
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
	carta := mazo[0]
	mazoRestante := mazo[1:]
	return carta, mazoRestante
}

func CalcularValor(cartas []string) int {
	total := 0
	ases := 0

	for _, carta := range cartas {
		// Normalizar la carta para asegurar que "a" o " A " se interpreten como "A"
		normalizedCarta := strings.ToUpper(strings.TrimSpace(carta))
		switch normalizedCarta {
		case "J", "Q", "K":
			total += 10
		case "A":
			total += 11
			ases++
		default:
			var valor int
			fmt.Sscanf(normalizedCarta, "%d", &valor) 
			total += valor
		}
	}

	// Ajustar Ases si el total es > 21 (cambiar 11 por 1)
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

	// Un Blackjack natural es un As y una carta de valor 10 (10, J, Q, K)
	valor1 := CalcularValor([]string{cartas[0]})
	valor2 := CalcularValor([]string{cartas[1]})

	// Determinar si alguna de las cartas es un As (valor 11 inicialmente)
	esAs1 := strings.ToUpper(strings.TrimSpace(cartas[0])) == "A"
	esAs2 := strings.ToUpper(strings.TrimSpace(cartas[1])) == "A"

	// Determinar si alguna de las cartas es un 10 (10, J, Q, K)
	esDiez1 := (valor1 == 10 && !esAs1) 
	esDiez2 := (valor2 == 10 && !esAs2)

	// Es Blackjack si uno es As y el otro es 10 (y la suma es 21)
	return (esAs1 && esDiez2) || (esAs2 && esDiez1)
}

func Es17Blando(cartas []string) bool {
	total := 0
	ases := 0

	for _, carta := range cartas {
		normalizedCarta := strings.ToUpper(strings.TrimSpace(carta))
		switch normalizedCarta {
		case "A":
			total += 11
			ases++
		case "J", "Q", "K":
			total += 10
		default:
			var valor int
			fmt.Sscanf(normalizedCarta, "%d", &valor)
			total += valor
		}
	}

	// Es un 17 blando si el total es 17 y al menos un As todavía se está contando como 11.
	return total == 17 && ases > 0 && (total-10 >= 1 && total-10 != 17) 
}

// EvaluarResultado determina el resultado de una mano de Blackjack.
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