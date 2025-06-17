package blackjack

import (
    "fmt"
    "math/rand"
    "strings"
)

// Mazo completo de 52 cartas (4 palos)
var mazoCompleto = []string{
    "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
    "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
    "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
    "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
}

// MezclarMazo devuelve un mazo aleatorio
func MezclarMazo() []string {
    baraja := append([]string{}, mazoCompleto...)
    rand.Shuffle(len(baraja), func(i, j int) {
        baraja[i], baraja[j] = baraja[j], baraja[i]
    })
    return baraja
}

// TomarCarta saca la primera carta del mazo y devuelve la carta y el mazo restante
func TomarCarta(mazo []string) (string, []string) {
    if len(mazo) == 0 {
        return "", mazo
    }
    return mazo[0], mazo[1:]
}

// CalcularValor evalúa el valor total de una mano
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

    // Ajuste de ases si se pasa de 21
    for total > 21 && ases > 0 {
        total -= 10
        ases--
    }

    return total
}

// CartasToString convierte []string a "A-10-Q"
func CartasToString(cartas []string) string {
    return strings.Join(cartas, "-")
}

// StringToCartas convierte "A-10-Q" a []string
func StringToCartas(cadena string) []string {
    if cadena == "" {
        return []string{}
    }
    return strings.Split(cadena, "-")
}
