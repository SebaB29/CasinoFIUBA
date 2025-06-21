package blackjack

import (
	"strings"
	"casino/models"
	repo "casino/repositories/juegos"
	"fmt"
	"github.com/gin-gonic/gin"
)

func finalizarPartida(partida *models.PartidaBlackjack, mazo []string) error {
	cartasBanca := StringToCartas(partida.CartasBanca)
	for CalcularValor(cartasBanca) < 17 && len(mazo) > 0 {
		var carta string
		carta, mazo = TomarCarta(mazo)
		cartasBanca = append(cartasBanca, carta)
	}
	partida.CartasBanca = CartasToString(cartasBanca)
	partida.Mazo = CartasToString(mazo)

	valBanca := CalcularValor(cartasBanca)
	cartas1 := StringToCartas(partida.CartasJugador)
	val1 := CalcularValor(cartas1)
	resultado1 := EvaluarResultado(val1, valBanca)

	if partida.CartasJugadorSplit != "" {
		cartas2 := StringToCartas(partida.CartasJugadorSplit)
		val2 := CalcularValor(cartas2)
		resultado2 := EvaluarResultado(val2, valBanca)

		switch {
		case resultado1 == "ganada" && resultado2 == "ganada":
			partida.Estado = models.Ganada
		case resultado1 == "perdida" && resultado2 == "perdida":
			partida.Estado = models.Perdida
		case resultado1 == "empatada" && resultado2 == "empatada":
			partida.Estado = models.Empatada
		default:
			partida.Estado = models.Perdida
		}
	} else {
		partida.Estado = EstadoFromResultado(resultado1)
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

	var cartas []string
	if partida.ManoActual == 2 {
		cartas = StringToCartas(partida.CartasJugadorSplit)
	} else {
		cartas = StringToCartas(partida.CartasJugador)
	}

	nuevaCarta, mazo := TomarCarta(mazo)
	cartas = append(cartas, nuevaCarta)
	valor := CalcularValor(cartas)

	// Guardamos la mano actualizada
	if partida.ManoActual == 2 {
		partida.CartasJugadorSplit = CartasToString(cartas)
	} else {
		partida.CartasJugador = CartasToString(cartas)
	}
	partida.Mazo = CartasToString(mazo)

	// Evaluamos si se pasó
	if valor > 21 {
		if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
			partida.ManoActual = 2
		} else {
			partida.Estado = models.Perdida
		}
		// Si es la segunda mano, finalizamos la partida completa
		if partida.ManoActual == 2 || partida.CartasJugadorSplit == "" {
			err = finalizarPartida(partida, mazo)
			if err != nil {
				return nil, fmt.Errorf("error al finalizar la partida: %w", err)
			}
			return gin.H{
				"estado": partida.Estado,
			}, nil
		}
	}

	// ✅ Si llegó a 21 y no hay segunda mano, finalizar
	if valor == 21 && (partida.CartasJugadorSplit == "" || partida.ManoActual == 2) {
		err = finalizarPartida(partida, mazo)
		if err != nil {
			return nil, fmt.Errorf("error al finalizar la partida: %w", err)
		}
		return gin.H{
			"estado": partida.Estado,
		}, nil
	}

	// Guardamos el estado si todavía no se cerró
	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("error actualizando partida: %w", err)
	}

	return gin.H{
		"mano":        partida.ManoActual,
		"cartas_mano": cartas,
		"valor":       valor,
		"estado":      partida.Estado,
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

	if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
		partida.ManoActual = 2
		if err := repo.ActualizarPartida(partida); err != nil {
			return nil, fmt.Errorf("error al pasar a la segunda mano: %w", err)
		}
		return gin.H{
			"mensaje": "Pasaste a la segunda mano",
			"mano":    2,
		}, nil
	}

	err = finalizarPartida(partida, mazo)
	if err != nil {
		return nil, fmt.Errorf("error finalizando partida: %w", err)
	}

	return gin.H{
		"cartas_banca":  StringToCartas(partida.CartasBanca),
		"cartas_mano_1": StringToCartas(partida.CartasJugador),
		"cartas_mano_2": StringToCartas(partida.CartasJugadorSplit),
		"estado":        partida.Estado,
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
	if partida.Doblar {
		return nil, fmt.Errorf("ya se ha doblado en esta partida")
	}

	mazo := StringToCartas(partida.Mazo)

	var cartas []string
	if partida.ManoActual == 2 {
		cartas = StringToCartas(partida.CartasJugadorSplit)
	} else {
		cartas = StringToCartas(partida.CartasJugador)
	}

	nuevaCarta, mazo := TomarCarta(mazo)
	cartas = append(cartas, nuevaCarta)
	valor := CalcularValor(cartas)

	partida.Apuesta *= 2
	partida.Doblar = true
	partida.Mazo = CartasToString(mazo)

	if partida.ManoActual == 2 {
		partida.CartasJugadorSplit = CartasToString(cartas)
	} else {
		partida.CartasJugador = CartasToString(cartas)
	}

	if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
		partida.ManoActual = 2
	} else {
		err = finalizarPartida(partida, mazo)
		if err != nil {
			return nil, fmt.Errorf("error al finalizar la partida: %w", err)
		}
	}

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("error actualizando la partida: %w", err)
	}

	return gin.H{
		"carta_nueva":     nuevaCarta,
		"mano":            partida.ManoActual,
		"valor_mano":      valor,
		"doblado":         true,
		"estado":          partida.Estado,
		"cartas_banca":    StringToCartas(partida.CartasBanca),
		"cartas_mano_1":   StringToCartas(partida.CartasJugador),
		"cartas_mano_2":   StringToCartas(partida.CartasJugadorSplit),
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
	if len(cartas) > 2 {
		return nil, fmt.Errorf("ya no puedes rendirte en esta fase del juego")
	}

	partida.Estado = models.Rendida
	partida.Apuesta = partida.Apuesta / 2

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo actualizar la partida al rendirse: %w", err)
	}

	return gin.H{
		"mensaje":          "Te has rendido, recuperas el 50% de tu apuesta",
		"estado":           partida.Estado,
		"apuesta_restante": partida.Apuesta,
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
	if len(cartas) != 2 || cartas[0] != cartas[1] {
		return nil, fmt.Errorf("solo se puede dividir si las dos cartas iniciales son iguales")
	}

	mazo := StringToCartas(partida.Mazo)

	cartaExtra1, mazo := TomarCarta(mazo)
	cartaExtra2, mazo := TomarCarta(mazo)

	mano1 := []string{cartas[0], cartaExtra1}
	mano2 := []string{cartas[1], cartaExtra2}

	partida.CartasJugador = CartasToString(mano1)
	partida.CartasJugadorSplit = CartasToString(mano2)
	partida.Mazo = CartasToString(mazo)
	partida.ApuestaSplit = partida.Apuesta
	partida.ManoActual = 1

	// Evaluación automática: si la mano 1 ya se pasó, pasamos a la segunda mano
	if CalcularValor(mano1) > 21 {
		partida.ManoActual = 2
	}

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo dividir la partida: %w", err)
	}

	return gin.H{
		"mensaje":       "Split realizado. Comenzando con la mano 1.",
		"mano_actual":   partida.ManoActual,
		"cartas_mano_1": mano1,
		"cartas_mano_2": mano2,
		"apuesta_split": partida.ApuestaSplit,
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

	apuestaSeguro := partida.Apuesta / 2
	var mensaje string

	if len(cartasBanca) == 0 || strings.ToUpper(strings.TrimSpace(cartasBanca[0])) != "A" {
		mensaje = "No se puede usar seguro: la banca no tiene un As como primera carta. La partida continúa normalmente."
		return gin.H{
			"mensaje":          mensaje,
			"seguro_pagado":    0,
			"seguro_resultado": 0,
			"estado":           partida.Estado,
		}, nil
	}

	tieneBlackjack := TieneBlackjack(cartasBanca)

	if tieneBlackjack {
		partida.Seguro = apuestaSeguro * 2
		partida.Estado = models.Perdida
		mensaje = "Seguro activado. La banca tiene blackjack. Ganás 2:1 pero perdiste la partida."
	} else {
		partida.Seguro = 0
		mensaje = "Seguro activado. La banca no tiene blackjack. Perdés el seguro, pero la partida continúa."
	}

	if err := repo.ActualizarPartida(partida); err != nil {
		return nil, fmt.Errorf("no se pudo actualizar la partida con el seguro: %w", err)
	}

	return gin.H{
		"mensaje":          mensaje,
		"seguro_pagado":    apuestaSeguro,
		"seguro_resultado": partida.Seguro,
		"estado":           partida.Estado,
	}, nil
}

