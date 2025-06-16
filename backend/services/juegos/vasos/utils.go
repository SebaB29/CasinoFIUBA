package vasos

func CalcularPremio(apuesta float64) float64 {
	// Devuelve 3x si gana (por cada 1 apostado gana 2 mas)
	return apuesta * 3
}
