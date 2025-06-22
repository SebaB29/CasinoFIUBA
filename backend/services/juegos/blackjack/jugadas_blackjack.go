package blackjack


import (
    "strings"
    "casino/models"
    repo "casino/repositories/juegos"
    "fmt"
    "github.com/gin-gonic/gin"
)


// finalizarPartida devuelve resultados por mano y actualiza el estado global
func finalizarPartida(partida *models.PartidaBlackjack, mazo []string) (string, string, error) {
    cartasBanca := StringToCartas(partida.CartasBanca)
    for CalcularValor(cartasBanca) < 17 && len(mazo) > 0 {
        carta, resto := TomarCarta(mazo)
        cartasBanca = append(cartasBanca, carta)
        mazo = resto
    }
    partida.CartasBanca = CartasToString(cartasBanca)
    partida.Mazo = CartasToString(mazo)

    valBanca := CalcularValor(cartasBanca)
    cartasJ1 := StringToCartas(partida.CartasJugador)
    valJ1 := CalcularValor(cartasJ1)
    res1 := EvaluarResultado(valJ1, valBanca)

    res2 := ""
    if partida.CartasJugadorSplit != "" {
        cartasJ2 := StringToCartas(partida.CartasJugadorSplit)
        valJ2 := CalcularValor(cartasJ2)
        res2 = EvaluarResultado(valJ2, valBanca)

        switch {
        case res1 == "ganada" || res2 == "ganada":
            partida.Estado = models.Ganada
        case res1 == "perdida" && res2 == "perdida":
            partida.Estado = models.Perdida
        case res1 == "empatada" && res2 == "empatada":
            partida.Estado = models.Empatada
        case (res1 == "empatada" && res2 == "perdida") || (res1 == "perdida" && res2 == "empatada"):
            partida.Estado = models.Perdida
        case (res1 == "ganada" && res2 == "empatada") || (res1 == "empatada" && res2 == "ganada"):
            partida.Estado = models.Ganada
        default:
            partida.Estado = models.Empatada
        }
    } else {
        partida.Estado = EstadoFromResultado(res1)
    }

    return res1, res2, repo.ActualizarPartida(partida)
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
    var mano []string
    if partida.ManoActual == 2 {
        mano = StringToCartas(partida.CartasJugadorSplit)
    } else {
        mano = StringToCartas(partida.CartasJugador)
    }

    if (partida.ManoActual == 1 && partida.IsSplitAcesMano1 && len(mano) >= 2) ||
       (partida.ManoActual == 2 && partida.IsSplitAcesMano2 && len(mano) >= 2) {
        return nil, fmt.Errorf("no puedes pedir más cartas en una mano de Ases divididos")
    }
    if len(mazo) == 0 {
        return nil, fmt.Errorf("el mazo está vacío")
    }

    nuevaCarta, mazo := TomarCarta(mazo)
    mano = append(mano, nuevaCarta)
    valor := CalcularValor(mano)

    if partida.ManoActual == 2 {
        partida.CartasJugadorSplit = CartasToString(mano)
    } else {
        partida.CartasJugador = CartasToString(mano)
    }
    partida.Mazo = CartasToString(mazo)

    if valor >= 21 {
        if partida.ManoActual == 1 && partida.CartasJugadorSplit != "" {
            partida.ManoActual = 2
        } else {
            res1, res2, err := finalizarPartida(partida, mazo)
            if err != nil {
                return nil, fmt.Errorf("error finalizando: %w", err)
            }
            return gin.H{
                "carta_nueva":      nuevaCarta,
                "estado":           partida.Estado,
                "resultado_mano_1": res1,
                "resultado_mano_2": res2,
                "cartas_banca":     StringToCartas(partida.CartasBanca),
                "valor_banca":      CalcularValor(StringToCartas(partida.CartasBanca)),
            }, nil
        }
    }

    if err := repo.ActualizarPartida(partida); err != nil {
        return nil, fmt.Errorf("error actualizando: %w", err)
    }
    return gin.H{
        "carta_nueva":         nuevaCarta,
        "mano_actual":         partida.ManoActual,
        "cartas_mano_1":       StringToCartas(partida.CartasJugador),
        "valor_mano_1":        CalcularValor(StringToCartas(partida.CartasJugador)),
        "cartas_mano_2":       StringToCartas(partida.CartasJugadorSplit),
        "valor_mano_2":        CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
        "cartas_banca":        []string{StringToCartas(partida.CartasBanca)[0], "???"},
        "valor_banca_visible": CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]}),
        "estado":              partida.Estado,
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
        return gin.H{
            "mensaje":             "Pasando a la segunda mano.",
            "mano_actual":         2,
            "cartas_mano_1":       StringToCartas(partida.CartasJugador),
            "valor_mano_1":        CalcularValor(StringToCartas(partida.CartasJugador)),
            "cartas_mano_2":       StringToCartas(partida.CartasJugadorSplit),
            "valor_mano_2":        CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
            "cartas_banca":        []string{StringToCartas(partida.CartasBanca)[0], "???"},
            "valor_banca_visible": CalcularValor([]string{StringToCartas(partida.CartasBanca)[0]}),
        }, nil
    }


    res1, res2, err := finalizarPartida(partida, mazo)
    if err != nil {
        return nil, fmt.Errorf("error finalizando: %w", err)
    }
    return gin.H{
        "cartas_banca":      StringToCartas(partida.CartasBanca),
        "valor_banca":       CalcularValor(StringToCartas(partida.CartasBanca)),
        "cartas_mano_1":     StringToCartas(partida.CartasJugador),
        "valor_mano_1":      CalcularValor(StringToCartas(partida.CartasJugador)),
        "cartas_mano_2":     StringToCartas(partida.CartasJugadorSplit),
        "valor_mano_2":      CalcularValor(StringToCartas(partida.CartasJugadorSplit)),
        "resultado_mano_1":  res1,
        "resultado_mano_2":  res2,
        "estado":            partida.Estado,
        "mensaje":           "Partida finalizada por Stand",
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
        return nil, fmt.Errorf("ya se ha doblado")
    }


    var mano []string
    if partida.ManoActual == 2 {
        mano = StringToCartas(partida.CartasJugadorSplit)
    } else {
        mano = StringToCartas(partida.CartasJugador)
    }
    if len(mano) != 2 {
        return nil, fmt.Errorf("solo puedes doblar con dos cartas iniciales")
    }


    mazo := StringToCartas(partida.Mazo)
    if len(mazo) == 0 {
        return nil, fmt.Errorf("mazo vacío")
    }


    nuevaCarta, mazo := TomarCarta(mazo)
    mano = append(mano, nuevaCarta)
    valor := CalcularValor(mano)


    if partida.ManoActual == 1 {
        partida.Apuesta *= 2
    } else {
        partida.ApuestaSplit *= 2
    }
    partida.Doblar = true


    if partida.ManoActual == 2 {
        partida.CartasJugadorSplit = CartasToString(mano)
    } else {
        partida.CartasJugador = CartasToString(mano)
    }
    partida.Mazo = CartasToString(mazo)


    res1, res2, err := finalizarPartida(partida, mazo)
    if err != nil {
        return nil, fmt.Errorf("error finalizando: %w", err)
    }


    return gin.H{
        "carta_nueva":        nuevaCarta,
        "estado":             partida.Estado,
        "resultado_mano_1":   res1,
        "resultado_mano_2":   res2,
        "apuesta_principal":  partida.Apuesta,
        "apuesta_split":      partida.ApuestaSplit,
        "cartas_banca":       StringToCartas(partida.CartasBanca),
        "valor_banca":        CalcularValor(StringToCartas(partida.CartasBanca)),
         "valor_mano_actual":  valor,
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
