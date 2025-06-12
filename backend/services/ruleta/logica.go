package services

import (
	dto "casino/dto/juegos"
	"math/rand"
	"time"
)

var colores = map[int]string{
	0: "verde",
	1: "rojo", 2: "negro", 3: "rojo", 4: "negro", 5: "rojo", 6: "negro",
	7: "rojo", 8: "negro", 9: "rojo", 10: "negro", 11: "negro", 12: "rojo",
	13: "negro", 14: "rojo", 15: "negro", 16: "rojo", 17: "negro", 18: "rojo",
	19: "rojo", 20: "negro", 21: "rojo", 22: "negro", 23: "rojo", 24: "negro",
	25: "rojo", 26: "negro", 27: "rojo", 28: "negro", 29: "negro", 30: "rojo",
	31: "negro", 32: "rojo", 33: "negro", 34: "rojo", 35: "negro", 36: "rojo",
}

func EjecutarRuleta(jugada dto.RuletaRequestDTO) dto.RuletaResponseDTO {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numeroGanador := r.Intn(37)
	colorGanador := colores[numeroGanador]

	multiplicador := calcularMultiplicador(jugada, numeroGanador, colorGanador)
	ganancia := jugada.Monto * multiplicador

	return dto.RuletaResponseDTO{
		MontoApostado: jugada.Monto,
		TipoApuesta:   jugada.TipoApuesta,
		Numeros:       jugada.Numeros,
		Docena:        jugada.Docena,
		Color:         jugada.Color,
		Paridad:       jugada.Paridad,
		AltoBajo:      jugada.AltoBajo,
		NumeroGanador: numeroGanador,
		ColorGanador:  colorGanador,
		Multiplicador: multiplicador,
		Ganancia:      ganancia,
	}
}

func calcularMultiplicador(jugada dto.RuletaRequestDTO, ganador int, colorGanador string) float64 {
	switch jugada.TipoApuesta {
	case "pleno":
		if contiene(jugada.Numeros, ganador) {
			return 36.0
		}
	case "dividida":
		if contiene(jugada.Numeros, ganador) {
			return 18.0
		}
	case "calle":
		if contiene(jugada.Numeros, ganador) {
			return 12.0
		}
	case "cuadro":
		if contiene(jugada.Numeros, ganador) {
			return 9.0
		}
	case "docena":
		if (jugada.Docena == 1 && ganador >= 1 && ganador <= 12) ||
			(jugada.Docena == 2 && ganador >= 13 && ganador <= 24) ||
			(jugada.Docena == 3 && ganador >= 25 && ganador <= 36) {
			return 3.0
		}
	case "color":
		if jugada.Color == colorGanador {
			return 2.0
		}
	case "paridad":
		if ganador != 0 && ((jugada.Paridad == "par" && ganador%2 == 0) || (jugada.Paridad == "impar" && ganador%2 != 0)) {
			return 2.0
		}
	case "alto_bajo":
		if (jugada.AltoBajo == "alto" && ganador >= 19 && ganador <= 36) ||
			(jugada.AltoBajo == "bajo" && ganador >= 1 && ganador <= 18) {
			return 2.0
		}
	}
	return 0.0
}

func contiene(numeros []int, objetivo int) bool {
	for _, n := range numeros {
		if n == objetivo {
			return true
		}
	}
	return false
}
