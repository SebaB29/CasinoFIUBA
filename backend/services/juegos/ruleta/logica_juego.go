package ruleta

import (
	dto "casino/dto/juegos"
	"math/rand"
	"time"
)

func obtenerNumeroGanador() NumeroRuleta {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numeroGanador := r.Intn(CantidadNumerosRuleta)

	if numeroGanador == 0 {
		return numeroCero
	}

	for _, fila := range tableroRuleta {
		for _, numero := range fila {
			if numero.Valor == numeroGanador {
				return numero
			}
		}
	}

	return NumeroRuleta{Valor: -1, Color: "desconocido"}
}

func calcularMultiplicador(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	switch jugada.TipoApuesta {
	case "pleno":
		if contiene(jugada.Numeros, numeroGanador.Valor) {
			return MultiplicadorPleno
		}
	case "dividida":
		if contiene(jugada.Numeros, numeroGanador.Valor) {
			return MultiplicadorDividida
		}
	case "calle":
		if contiene(jugada.Numeros, numeroGanador.Valor) {
			return MultiplicadorCalle
		}
	case "cuadro":
		if contiene(jugada.Numeros, numeroGanador.Valor) {
			return MultiplicadorCuadro
		}
	case "docena":
		if (jugada.Docena == PrimeraDocena && numeroGanador.Valor >= MinDocena1 && numeroGanador.Valor <= MaxDocena1) ||
			(jugada.Docena == SegundaDocena && numeroGanador.Valor >= MinDocena2 && numeroGanador.Valor <= MaxDocena2) ||
			(jugada.Docena == TerceraDocena && numeroGanador.Valor >= MinDocena3 && numeroGanador.Valor <= MaxDocena3) {
			return MultiplicadorDocena
		}
	case "color":
		if jugada.Color == numeroGanador.Color {
			return MultiplicadorSimple
		}
	case "paridad":
		if numeroGanador.Valor != 0 && ((jugada.Paridad == "par" && numeroGanador.Valor%2 == 0) || (jugada.Paridad == "impar" && numeroGanador.Valor%2 != 0)) {
			return MultiplicadorSimple
		}
	case "alto_bajo":
		if (jugada.AltoBajo == "alto" && numeroGanador.Valor >= MinAlto && numeroGanador.Valor <= MaxAlto) ||
			(jugada.AltoBajo == "bajo" && numeroGanador.Valor >= MinBajo && numeroGanador.Valor <= MaxBajo) {
			return MultiplicadorSimple
		}
	}
	return SinGanancia
}

func contiene(numeros []int, objetivo int) bool {
	for _, numero := range numeros {
		if numero == objetivo {
			return true
		}
	}
	return false
}