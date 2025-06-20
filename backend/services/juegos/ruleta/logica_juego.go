package ruleta

import (
	dto "casino/dto/juegos"
	"math/rand"
	"time"
)

var calculadores = map[string]func(dto.RuletaRequestDTO, NumeroRuleta) float64{
	"pleno":     calcularPleno,
	"dividida":  calcularDividida,
	"calle":     calcularCalle,
	"cuadro":    calcularCuadro,
	"docena":    calcularDocena,
	"color":     calcularColor,
	"paridad":   calcularParidad,
	"alto_bajo": calcularAltoBajo,
}

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
	if fn, ok := calculadores[jugada.TipoApuesta]; ok {
		return fn(jugada, numeroGanador)
	}
	return SinGanancia
}

func extraerDetalles(jugada dto.RuletaRequestDTO) interface{} {
	switch jugada.TipoApuesta {
	case "pleno", "dividida", "calle", "cuadro":
		return map[string]interface{}{"numeros": jugada.Numeros}
	case "docena":
		return map[string]interface{}{"docena": jugada.Docena}
	case "color":
		return map[string]interface{}{"color": jugada.Color}
	case "paridad":
		return map[string]interface{}{"paridad": jugada.Paridad}
	case "alto_bajo":
		return map[string]interface{}{"alto_bajo": jugada.AltoBajo}
	default:
		return nil
	}
}

// ----------- Calculadores -----------

func calcularPleno(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if contiene(jugada.Numeros, numeroGanador.Valor) {
		return MultiplicadorPleno
	}
	return SinGanancia
}

func calcularDividida(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if contiene(jugada.Numeros, numeroGanador.Valor) {
		return MultiplicadorDividida
	}
	return SinGanancia
}

func calcularCalle(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if contiene(jugada.Numeros, numeroGanador.Valor) {
		return MultiplicadorCalle
	}
	return SinGanancia
}

func calcularCuadro(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if contiene(jugada.Numeros, numeroGanador.Valor) {
		return MultiplicadorCuadro
	}
	return SinGanancia
}

func calcularDocena(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	switch jugada.Docena {
	case PrimeraDocena:
		if numeroGanador.Valor >= MinDocena1 && numeroGanador.Valor <= MaxDocena1 {
			return MultiplicadorDocena
		}
	case SegundaDocena:
		if numeroGanador.Valor >= MinDocena2 && numeroGanador.Valor <= MaxDocena2 {
			return MultiplicadorDocena
		}
	case TerceraDocena:
		if numeroGanador.Valor >= MinDocena3 && numeroGanador.Valor <= MaxDocena3 {
			return MultiplicadorDocena
		}
	}
	return SinGanancia
}

func calcularColor(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if jugada.Color == numeroGanador.Color {
		return MultiplicadorSimple
	}
	return SinGanancia
}

func calcularParidad(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if numeroGanador.Valor == 0 {
		return SinGanancia
	}
	if (jugada.Paridad == "par" && numeroGanador.Valor%2 == 0) || (jugada.Paridad == "impar" && numeroGanador.Valor%2 != 0) {
		return MultiplicadorSimple
	}
	return SinGanancia
}

func calcularAltoBajo(jugada dto.RuletaRequestDTO, numeroGanador NumeroRuleta) float64 {
	if jugada.AltoBajo == "alto" && numeroGanador.Valor >= MinAlto && numeroGanador.Valor <= MaxAlto {
		return MultiplicadorSimple
	}
	if jugada.AltoBajo == "bajo" && numeroGanador.Valor >= MinBajo && numeroGanador.Valor <= MaxBajo {
		return MultiplicadorSimple
	}
	return SinGanancia
}

// ----------- Funciones Auxiliares -----------

func contiene(numeros []int, objetivo int) bool {
	for _, numero := range numeros {
		if numero == objetivo {
			return true
		}
	}
	return false
}
