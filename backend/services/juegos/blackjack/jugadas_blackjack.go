package blackjack

import (
	"strings"
	"casino/models"
	repo "casino/repositories/juegos"
	"fmt"
	"github.com/gin-gonic/gin"
)

// finalizarPartida ahora se encargará de establecer el 'Estado' final de la partida
func finalizarPartida(partida *models.PartidaBlackjack, mazo []string) error {
	cartasBanca := StringToCartas(partida.CartasBanca)

	// La banca toma cartas hasta alcanzar 17 o más, o quedarse sin mazo.
	for CalcularValor(cartasBanca) < 17 && len(mazo) > 0 {
		var carta string
		carta, mazo = TomarCarta(mazo)
		cartasBanca = append(cartasBanca, carta)
	}
	partida.CartasBanca = CartasToString(cartasBanca)
	partida.Mazo = CartasToString(mazo)

	valBanca := CalcularValor(cartasBanca)
	cartasJugador1 := StringToCartas(partida.CartasJugador)
	valJugador1 := CalcularValor(cartasJugador1)

	// Evaluar resultado para la mano principal
	resultadoMano1 := EvaluarResultado(valJugador1, valBanca)

	// Si hay split, evaluar también la segunda mano
	if partida.CartasJugadorSplit != "" {
		cartasJugadorSplit := StringToCartas(partida.CartasJugadorSplit)
		valJugadorSplit := CalcularValor(cartasJugadorSplit)
		resultadoMano2 := EvaluarResultado(valJugadorSplit, valBanca)

		switch {
		case resultadoMano1 == "ganada" || resultadoMano2 == "ganada":
			partida.Estado = models.Ganada
		case resultadoMano1 == "perdida" && resultadoMano2 == "perdida":
			partida.Estado = models.Perdida
		case resultadoMano1 == "empatada" && resultadoMano2 == "empatada":
			partida.Estado = models.Empatada
		case resultadoMano1 == "empatada" && resultadoMano2 == "perdida" || resultadoMano1 == "perdida" && resultadoMano2 == "empatada":
			partida.Estado = models.Perdida 
		case resultadoMano1 == "ganada" && resultadoMano2 == "empatada" || resultadoMano1 == "empatada" && resultadoMano2 == "ganada":
			partida.Estado = models.Ganada 
		default:
			partida.Estado = models.Empatada 
		}
	} else {
		// Si no hay split, el estado global es el resultado de la única mano
		partida.Estado = EstadoFromResultado(resultadoMano1)
	}

	return repo.ActualizarPartida(partida)
}

func (s *BlackjackService) Hit(idPartida uint) (gin.H, error) {
	partida, err := repo.ObtenerPartidaPorID(idPartida)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}
	if partida.Estado != models.EnCurso {
		return nil, fmt.Errorf("la partida ya está finalizada")
	}

	mazo := StringToCartas(partida.Mazo)

	var cartasManoActual []string
	if partida.ManoActual == 2 {
		cartasManoActual = StringToCartas(partida.CartasJugadorSplit)
	} else {
		cartasManoActual = StringToCartas(partida.CartasJugador)
	}

	 if (partida.ManoActual == 1 && partida.IsSplitAcesMano1 && len(cartasManoActual) >= 2) ||
	    (partida.ManoActual == 2 && partida.IsSplitAcesMano2 && len(cartasManoActual) >= 2) {
	    return nil, fmt.Errorf("no puedes pedir más cartas en una mano de Ases divididos.")
	 }

	if len(mazo) == 0 {
		return nil, fmt.Errorf("el mazo está vacío, no se pueden tomar más cartas")
	}

	nuevaCarta, mazo := TomarCarta(mazo)
	cartasManoActual = append(cartasManoActual, nuevaCarta)
	valorManoActual := CalcularValor(cartasManoActual)

	// Guardamos la mano actualizada
	if partida.ManoActual == 2 {
		partida.CartasJugadorSplit = CartasToString(cartasManoActual)
	} else {
		partida.CartasJugador = CartasToString(cartasManoActual)
	}
	partida.Mazo = CartasToString(mazo)

	// Evaluamos si la mano actual se pasó (Bust)
	if valorManoActual > 21 {
		if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
			// Si la primera mano se pasa y hay una segunda, pasamos a la segunda mano
			partida.ManoActual = 2
		} else {
			// Si es la segunda mano o no hay split, la partida finaliza
			err = finalizarPartida(partida, mazo)
			if err != nil {
				return nil, fmt.Errorf("error al finalizar la partida por bust: %w", err)
			}
		}
	} else if valorManoActual == 21 {
		// Si la mano actual llega a 21
		if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
			// Si es la primera mano y hay una segunda, pasamos a la segunda mano
			partida.ManoActual = 2
		} else {
			// Si es la segunda mano o no hay split, la partida finaliza
			err = finalizarPartida(partida, mazo)
			if err != nil {
				return nil, fmt.Errorf("error al finalizar la partida por 21: %w", err)
			}
		}
	}

	// Guardamos el estado si todavía no se cerró
	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("error actualizando partida: %w", err)
	}

	// Preparamos la respuesta para el cliente, mostrando solo la primera carta de la banca
	cartasBancaVisible := []string{}
	valorBancaVisible := 0
	if len(StringToCartas(partida.CartasBanca)) > 0 {
		cartasBancaVisible = []string{StringToCartas(partida.CartasBanca)[0], "???"}
		valorBancaVisible = CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]})
	}

	return gin.H{
		"mano_actual":         partida.ManoActual,
		"cartas_mano_1":       StringToCartas(partida.CartasJugador),
		"valor_mano_1":        CalcularValor(StringToCartas(partida.CartasJugador)),
		"cartas_mano_2":       StringToCartas(partida.CartasJugadorSplit),
		"valor_mano_2":        CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
		"cartas_banca":        cartasBancaVisible,
		"valor_banca_visible": valorBancaVisible,
		"estado":              partida.Estado,
		"carta_nueva":         nuevaCarta, 
	}, nil
}

func (s *BlackjackService) Stand(idPartida uint) (gin.H, error) {
	partida, err := repo.ObtenerPartidaPorID(idPartida)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}
	if partida.Estado != models.EnCurso {
		return nil, fmt.Errorf("la partida ya fue finalizada")
	}

	mazo := StringToCartas(partida.Mazo)

	// Si hay una segunda mano y el jugador se planta en la primera, cambia a la segunda mano
	if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
		partida.ManoActual = 2
		if err := repo.ActualizarPartida(partida); err != nil {
			return nil, fmt.Errorf("error al pasar a la segunda mano: %w", err)
		}
		// Preparamos la respuesta para el cliente, mostrando solo la primera carta de la banca
		cartasBancaVisible := []string{}
		valorBancaVisible := 0
		if len(StringToCartas(partida.CartasBanca)) > 0 {
			cartasBancaVisible = []string{StringToCartas(partida.CartasBanca)[0], "???"}
			valorBancaVisible = CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]})
		}
		return gin.H{
			"mensaje":             "Te plantaste en la primera mano. Pasando a la segunda mano.",
			"mano_actual":         2,
			"estado":              partida.Estado,
			"cartas_mano_1":       StringToCartas(partida.CartasJugador),
			"valor_mano_1":        CalcularValor(StringToCartas(partida.CartasJugador)),
			"cartas_mano_2":       StringToCartas(partida.CartasJugadorSplit),
			"valor_mano_2":        CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
			"cartas_banca":        cartasBancaVisible,
			"valor_banca_visible": valorBancaVisible,
		}, nil
	}

	// Si no hay segunda mano o ya se está en la segunda, finalizar la partida
	err = finalizarPartida(partida, mazo)
	if err != nil {
		return nil, fmt.Errorf("error finalizando partida: %w", err)
	}

	// Devolvemos todos los detalles finales para que el frontend pueda mostrar el resultado
	return gin.H{
		"cartas_banca":    StringToCartas(partida.CartasBanca),
		"valor_banca":     CalcularValor(StringToCartas(partida.CartasBanca)),
		"cartas_mano_1":   StringToCartas(partida.CartasJugador),
		"valor_mano_1":    CalcularValor(StringToCartas(partida.CartasJugador)),
		"cartas_mano_2":   StringToCartas(partida.CartasJugadorSplit),
		"valor_mano_2":    CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
		"estado":          partida.Estado, // Este será el estado global final
		"mensaje":         "Partida finalizada por Stand.",
	}, nil
}

func (s *BlackjackService) Doblar(idPartida uint) (gin.H, error) {
	partida, err := repo.ObtenerPartidaPorID(idPartida)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}
	if partida.Estado != models.EnCurso {
		return nil, fmt.Errorf("la partida ya fue finalizada")
	}
	// `partida.Doblar` debe ser por partida, no por mano si solo hay un campo
	if partida.Doblar {
		return nil, fmt.Errorf("ya se ha doblado en esta partida.")
	}

	var cartasManoActual []string
	if partida.ManoActual == 2 {
		cartasManoActual = StringToCartas(partida.CartasJugadorSplit)
	} else {
		cartasManoActual = StringToCartas(partida.CartasJugador)
	}

	// Validación clave: Solo se puede doblar con las dos cartas iniciales de la mano actual.
	if len(cartasManoActual) != 2 {
		return nil, fmt.Errorf("solo puedes doblar con tus dos primeras cartas de la mano actual.")
	}

	mazo := StringToCartas(partida.Mazo)
	if len(mazo) == 0 {
		return nil, fmt.Errorf("el mazo está vacío, no se puede tomar más cartas para doblar")
	}

	nuevaCarta, mazo := TomarCarta(mazo)
	cartasManoActual = append(cartasManoActual, nuevaCarta)
	valorManoActual := CalcularValor(cartasManoActual)

	// Duplica la apuesta principal si la acción de doblar es para la mano principal.
	if partida.ManoActual == 1 {
		partida.Apuesta *= 2
	} else { 
		partida.ApuestaSplit *= 2
	}
	partida.Doblar = true 

	// Guardamos la mano actualizada en el modelo
	if partida.ManoActual == 2 {
		partida.CartasJugadorSplit = CartasToString(cartasManoActual)
	} else {
		partida.CartasJugador = CartasToString(cartasManoActual)
	}
	partida.Mazo = CartasToString(mazo)

	// Si hay split y se dobló la primera mano, se pasa a la segunda mano.
	// Si se dobló la segunda mano o no había split, la partida finaliza.
	if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
		partida.ManoActual = 2
	} else {
		err = finalizarPartida(partida, mazo)
		if err != nil {
			return nil, fmt.Errorf("error al finalizar la partida después de doblar: %w", err)
		}
	}

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("error actualizando la partida después de doblar: %w", err)
	}

	// Preparamos la respuesta para el cliente, mostrando solo la primera carta de la banca
	cartasBancaVisible := []string{}
	valorBancaVisible := 0
	if len(StringToCartas(partida.CartasBanca)) > 0 {
		cartasBancaVisible = []string{StringToCartas(partida.CartasBanca)[0], "???"}
		valorBancaVisible = CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]})
	}

	return gin.H{
		"carta_nueva":         nuevaCarta,
		"mano_actual":         partida.ManoActual,
		"valor_mano_actual":   valorManoActual, 
		"doblado":             true,
		"estado":              partida.Estado,
		"apuesta_principal":   partida.Apuesta,
		"apuesta_split":       partida.ApuestaSplit,
		"cartas_mano_1":       StringToCartas(partida.CartasJugador),
		"valor_mano_1":        CalcularValor(StringToCartas(partida.CartasJugador)),
		"cartas_mano_2":       StringToCartas(partida.CartasJugadorSplit),
		"valor_mano_2":        CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
		"cartas_banca":        cartasBancaVisible,
		"valor_banca_visible": valorBancaVisible,
	}, nil
}

func (s *BlackjackService) Rendirse(idPartida uint) (gin.H, error) {
	partida, err := repo.ObtenerPartidaPorID(idPartida)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}
	if partida.Estado != models.EnCurso {
		return nil, fmt.Errorf("la partida ya fue finalizada")
	}

	cartas := StringToCartas(partida.CartasJugador)
	// Solo puede rendirse si no ha pedido cartas adicionales (es decir, tiene 2 cartas iniciales)
	if len(cartas) > 2 {
		return nil, fmt.Errorf("no puedes rendirte en esta fase del juego (solo es posible con las dos primeras cartas).")
	}
	// Normalmente, rendirse no está permitido si se ha hecho split.
	if partida.CartasJugadorSplit != "" {
		return nil, fmt.Errorf("no puedes rendirte si has dividido tus cartas.")
	}

	partida.Estado = models.Rendida 
	partida.Apuesta = partida.Apuesta / 2 

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo actualizar la partida al rendirse: %w", err)
	}

	cartasBancaVisible := []string{}
	valorBancaVisible := 0
	if len(StringToCartas(partida.CartasBanca)) > 0 {
		cartasBancaVisible = []string{StringToCartas(partida.CartasBanca)[0], "???"}
		valorBancaVisible = CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]})
	}

	return gin.H{
		"mensaje":          "Te has rendido. Recuperas el 50% de tu apuesta.",
		"estado":           partida.Estado,
		"apuesta_restante": partida.Apuesta,
		"cartas_mano_1":    cartas,
		"valor_mano_1":     CalcularValor(cartas),
		"cartas_banca":     cartasBancaVisible,
		"valor_banca_visible": valorBancaVisible,
	}, nil
}

func (s *BlackjackService) Split(idPartida uint) (gin.H, error) {
	partida, err := repo.ObtenerPartidaPorID(idPartida)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}
	if partida.Estado != models.EnCurso {
		return nil, fmt.Errorf("la partida ya fue finalizada")
	}

	cartas := StringToCartas(partida.CartasJugador)
	// Validación clave: Solo se puede dividir si se tienen EXACTAMENTE dos cartas y son del mismo valor.
	if len(cartas) != 2 || CalcularValor([]string{cartas[0]}) != CalcularValor([]string{cartas[1]}) {
		return nil, fmt.Errorf("solo se puede dividir si las dos cartas iniciales son del mismo valor.")
	}
	// No se puede hacer split si ya hay una mano dividida (no se permite re-split en esta implementación)
	if partida.CartasJugadorSplit != "" {
		return nil, fmt.Errorf("ya has dividido una mano en esta partida.")
	}
	// No se puede doblar y hacer split
	if partida.Doblar {
		return nil, fmt.Errorf("no puedes dividir después de haber doblado.")
	}
	// No se puede rendir y hacer split (aunque la validación de rendirse ya cubre esto)
	if partida.Estado == models.Rendida {
		return nil, fmt.Errorf("no puedes dividir si ya te has rendido.")
	}

	mazo := StringToCartas(partida.Mazo)
	if len(mazo) < 2 {
		return nil, fmt.Errorf("no hay suficientes cartas en el mazo para realizar un split")
	}

	// Se reparte una carta a cada nueva mano
	cartaExtra1, mazo := TomarCarta(mazo)
	cartaExtra2, mazo := TomarCarta(mazo)

	mano1 := []string{cartas[0], cartaExtra1}
	mano2 := []string{cartas[1], cartaExtra2}

	partida.CartasJugador = CartasToString(mano1)
	partida.CartasJugadorSplit = CartasToString(mano2)
	partida.Mazo = CartasToString(mazo)
	partida.ApuestaSplit = partida.Apuesta 
	partida.ManoActual = 1 

	 if strings.ToUpper(strings.TrimSpace(cartas[0])) == "A" {
	    partida.IsSplitAcesMano1 = true 
	    partida.IsSplitAcesMano2 = true
	 }

	// Evaluación automática: si la mano 1 ya se pasó (bust) con la nueva carta, pasamos a la segunda mano
	if CalcularValor(mano1) > 21 {
		partida.ManoActual = 2
	}

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo dividir la partida: %w", err)
	}

	// Preparamos la respuesta para el cliente, mostrando solo la primera carta de la banca
	cartasBancaVisible := []string{}
	valorBancaVisible := 0
	if len(StringToCartas(partida.CartasBanca)) > 0 {
		cartasBancaVisible = []string{StringToCartas(partida.CartasBanca)[0], "???"}
		valorBancaVisible = CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]})
	}

	return gin.H{
		"mensaje":             "Split realizado. Comenzando con la mano 1.",
		"mano_actual":         partida.ManoActual,
		"cartas_mano_1":       mano1,
		"valor_mano_1":        CalcularValor(mano1),
		"cartas_mano_2":       mano2,
		"valor_mano_2":        CalcularValor(mano2),
		"apuesta_principal":   partida.Apuesta,
		"apuesta_split":       partida.ApuestaSplit,
		"estado":              partida.Estado,
		"cartas_banca":        cartasBancaVisible,
		"valor_banca_visible": valorBancaVisible,
	}, nil
}

func (s *BlackjackService) Seguro(idPartida uint) (gin.H, error) {
	partida, err := repo.ObtenerPartidaPorID(idPartida)
	if err != nil {
		return nil, fmt.Errorf("partida no encontrada")
	}
	if partida.Estado != models.EnCurso {
		return nil, fmt.Errorf("la partida ya fue finalizada")
	}

	cartasBanca := StringToCartas(partida.CartasBanca)

	// El seguro solo se puede tomar si la primera carta visible de la banca es un As.
	if len(StringToCartas(partida.CartasJugador)) > 2 || partida.CartasJugadorSplit != "" || partida.Doblar || partida.Seguro > 0 {
		return nil, fmt.Errorf("el seguro solo puede tomarse al inicio de la partida antes de cualquier otra acción.")
	}
	if len(cartasBanca) == 0 || strings.ToUpper(strings.TrimSpace(cartasBanca[0])) != "A" {
		mensaje := "No se puede usar seguro: la banca no tiene un As como primera carta. La partida continúa normalmente."
		return gin.H{
			"mensaje":          mensaje,
			"seguro_pagado":    0,
			"seguro_resultado": 0,
			"estado":           partida.Estado,
			"cartas_banca":     []string{cartasBanca[0], "???"}, 
		}, nil
	}

	apuestaSeguro := partida.Apuesta / 2

	tieneBlackjackLaBanca := TieneBlackjack(cartasBanca) 

	var mensaje string
	if tieneBlackjackLaBanca {
		partida.Seguro = apuestaSeguro * 2 
		partida.Estado = models.Perdida   
		mensaje = "Seguro activado. La banca tiene Blackjack. Ganás 2:1 en el seguro, pero perdiste la partida principal."
	} else {
		partida.Seguro = 0 
		mensaje = "Seguro activado. La banca no tiene Blackjack. Perdés el seguro, pero la partida principal continúa."
	}

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo actualizar la partida con el seguro: %w", err)
	}

	// Si la banca tenía Blackjack (y por ende se activó el seguro), se revelan las cartas de la banca.
	cartasBancaParaCliente := []string{}
	valorBancaParaCliente := 0
	if tieneBlackjackLaBanca {
		cartasBancaParaCliente = cartasBanca
		valorBancaParaCliente = CalcularValor(cartasBanca)
	} else {
		cartasBancaParaCliente = []string{cartasBanca[0], "???"}
		valorBancaParaCliente = CalcularValor([]string{cartasBanca[0]})
	}

	return gin.H{
		"mensaje":             mensaje,
		"seguro_pagado":       apuestaSeguro, 
		"seguro_resultado":    partida.Seguro, 
		"estado":              partida.Estado,
		"cartas_banca":        cartasBancaParaCliente,
		"valor_banca_visible": valorBancaParaCliente,
	}, nil
}