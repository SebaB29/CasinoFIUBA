package blackjack


import (
    "log"
    "strings"
    dto "casino/dto/juegos"
    "casino/models"
    repo "casino/repositories/juegos"
    "fmt"
)

// todosLosJugadoresActivosTerminaron verifica si todos los jugadores activos de la mesa han terminado su mano
func todosLosJugadoresActivosTerminaron(mesa *models.MesaBlackjack) bool {
    for _, mano := range mesa.ManosJugadores {
        if mano.Estado == string(models.EnCurso) {
            return false
        }
    }
    return true
}
// avanzarManoOMesa avanza a la siguiente mano del jugador o al siguiente jugador y finaliza la mesa si es necesario
func (s *BlackjackService) avanzarManoOMesa(mano *models.ManoJugadorBlackjack, mesa *models.MesaBlackjack) {
    // Si tiene split y falta jugar la segunda mano
    if mano.ManoActual == 1 && mano.CartasSplit != "" && mano.ResultadoMano2 == "" {
        mano.ManoActual = 2
        return
    }

    // actualizar el estado de la mano
    mano.Estado = CalcularEstadoFinal(mano.ResultadoMano1, mano.ResultadoMano2)

    // chequea si todos terminaron
    if todosLosJugadoresActivosTerminaron(mesa) {
        if err := finalizarMesa(mesa); err == nil {
            nuevaMesa, err := repo.ObtenerMesaConManos(mesa.ID)
            if err == nil {
                *mesa = *nuevaMesa
            }
        }
        return
    }

    // Si no todos terminaron, avanza al siguiente jugador
    s.AvanzarTurno(mesa)
}
// finalizarMesa resuelve la banca y determina el resultado para cada jugador
func finalizarMesa(mesa *models.MesaBlackjack) error {
    cartasBanca := StringToCartas(mesa.CartasBanca)
    mazo := StringToCartas(mesa.Mazo)

    // Banca roba hasta llegar a 17 o más
    for CalcularValor(cartasBanca) < 17 && len(mazo) > 0 {
        var carta string
        carta, mazo = TomarCarta(mazo)
        cartasBanca = append(cartasBanca, carta)
    }

    valorBanca := CalcularValor(cartasBanca)
    mesa.CartasBanca = CartasToString(cartasBanca)
    mesa.Mazo = CartasToString(mazo)

    // Evaluar cada mano de los jugadores
    for i := range mesa.ManosJugadores {
        mano := &mesa.ManosJugadores[i]

        // Saltear solo si ya se resolvió por completo (ambas manos si hay split)
        if mano.ResultadoMano1 != "" && (mano.CartasSplit == "" || mano.ResultadoMano2 != "") {
            mano.Estado = CalcularEstadoFinal(mano.ResultadoMano1, mano.ResultadoMano2)
            continue
        }

        // Mano principal
        if mano.ResultadoMano1 == "" {
            cartas := StringToCartas(mano.Cartas)
            valor1 := CalcularValor(cartas)
            mano.ResultadoMano1 = EvaluarResultado(valor1, valorBanca)
        }

        // Mano dividida
        if mano.CartasSplit != "" && mano.ResultadoMano2 == "" {
            cartasSplit := StringToCartas(mano.CartasSplit)
            valor2 := CalcularValor(cartasSplit)
            mano.ResultadoMano2 = EvaluarResultado(valor2, valorBanca)
        }

        // Siempre actualiza el estado final
        mano.Estado = CalcularEstadoFinal(mano.ResultadoMano1, mano.ResultadoMano2)
    }

    // Marca la mesa como finalizada
    mesa.Estado = "finalizada"

    if err := repo.ActualizarMesa(mesa); err != nil {
        return fmt.Errorf("fallo al finalizar la mesa: %w", err)
    }

    return nil
}

func (s *BlackjackService) Hit(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    if mesa.Estado != "en_curso" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mesa ya está finalizada")
    }

    indiceMano, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if mesa.JugadorActual != indiceMano {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es tu turno todavía")
    }

    mano := &mesa.ManosJugadores[indiceMano]
    if mano.Estado != string(models.EnCurso) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("esta mano ya está resuelta")
    }

    if (mano.ManoActual == 1 && mano.ResultadoMano1 != "") || (mano.ManoActual == 2 && mano.ResultadoMano2 != "") {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("esta mano ya fue jugada")
    }

    mazo := StringToCartas(mesa.Mazo)
    if len(mazo) == 0 {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no quedan cartas en el mazo")
    }

    carta, mazo := TomarCarta(mazo)

    var cartas []string
    if mano.ManoActual == 2 {
        cartas = StringToCartas(mano.CartasSplit)
    } else {
        cartas = StringToCartas(mano.Cartas)
    }

    cartas = append(cartas, carta)
    valor := CalcularValor(cartas)

    if mano.ManoActual == 2 {
        mano.CartasSplit = CartasToString(cartas)
    } else {
        mano.Cartas = CartasToString(cartas)
    }

    mesa.Mazo = CartasToString(mazo)

    // Evaluar si se pasa o llega a 21
    if valor >= 21 {
        if valor > 21 {
            if mano.ManoActual == 2 {
                mano.ResultadoMano2 = string(models.Perdida)
            } else {
                mano.ResultadoMano1 = string(models.Perdida)
            }
        } else {
            if mano.ManoActual == 2 {
                mano.ResultadoMano2 = string(models.Ganada)
            } else {
                mano.ResultadoMano1 = string(models.Ganada)
            }
        }

        mano.Estado = CalcularEstadoFinal(mano.ResultadoMano1, mano.ResultadoMano2)
        s.avanzarManoOMesa(mano, mesa)
    } else {
        // Mano sigue en juego, solo actualizamos la mesa
        if err := repo.ActualizarMesa(mesa); err != nil {
            return dto.BlackjackEstadoDTO{}, fmt.Errorf("error actualizando mesa: %w", err)
        }

        return s.EstadoParaJugador(mesa, userID), nil
    }

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("error actualizando mesa: %w", err)
    }

    log.Printf("Jugador actual post-hit: %d", mesa.JugadorActual)
    log.Printf("JugadorActual: %d, UserID solicitado: %d, IndiceJugador: %d", mesa.JugadorActual, userID, indiceMano)

    return s.EstadoParaJugador(mesa, userID), nil
}

func (s *BlackjackService) Stand(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    if mesa.Estado != "en_curso" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mesa ya ha finalizado")
    }

    indiceMano, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if mesa.JugadorActual != indiceMano {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es tu turno todavía")
    }

    mano := &mesa.ManosJugadores[indiceMano]

    if mano.Estado != string(models.EnCurso) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("esta mano ya fue resuelta")
    }

    // Avanzar a la siguiente mano o jugador
    s.avanzarManoOMesa(mano, mesa)

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("error actualizando mesa: %w", err)
    }

    return s.EstadoParaJugador(mesa, userID), nil
}

func (s *BlackjackService) Doblar(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    if mesa.Estado != "en_curso" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mesa ya ha finalizado")
    }

    indiceMano, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if mesa.JugadorActual != indiceMano {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es tu turno todavía")
    }

    mano := &mesa.ManosJugadores[indiceMano]

    if mano.Estado != string(models.EnCurso) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("esta mano ya está resuelta")
    }

    if mano.ManoActual != 1 || mano.CartasSplit != "" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("solo puedes doblar con una mano sin dividir")
    }

    cartas := StringToCartas(mano.Cartas)
    if len(cartas) != 2 {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("solo puedes doblar con dos cartas iniciales")
    }

    // Duplicar apuesta
    mano.Apuesta *= 2

    // Robar una carta
    mazo := StringToCartas(mesa.Mazo)
    if len(mazo) == 0 {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mazo vacío")
    }

    carta, mazo := TomarCarta(mazo)
    cartas = append(cartas, carta)
    mano.Cartas = CartasToString(cartas)
    mesa.Mazo = CartasToString(mazo)

    // Evaluar resultado automáticamente tras doblar
    valor := CalcularValor(cartas)
    if valor > 21 {
        mano.ResultadoMano1 = string(models.Perdida)
    }

    mano.Estado = CalcularEstadoFinal(mano.ResultadoMano1, mano.ResultadoMano2)

    s.avanzarManoOMesa(mano, mesa)

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("error actualizando mesa: %w", err)
    }

    return s.EstadoParaJugador(mesa, userID), nil
}

func (s *BlackjackService) Rendirse(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    if mesa.Estado != "en_curso" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mesa ya ha finalizado")
    }

    indiceMano, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if mesa.JugadorActual != indiceMano {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es tu turno todavía")
    }

    mano := &mesa.ManosJugadores[indiceMano]

    if mano.Estado != string(models.EnCurso) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("esta mano ya está resuelta")
    }

    cartas := StringToCartas(mano.Cartas)
    if len(cartas) != 2 || mano.CartasSplit != "" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("solo puedes rendirte con dos cartas sin haber dividido")
    }

    mano.Estado = string(models.Rendida)
    mano.Apuesta /= 2

    s.avanzarManoOMesa(mano, mesa)

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("error actualizando mesa: %w", err)
    }

    return s.EstadoParaJugador(mesa, userID), nil
}

func (s *BlackjackService) Split(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    if mesa.Estado != "en_curso" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mesa ya ha finalizado")
    }

    indiceMano, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if mesa.JugadorActual != indiceMano {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es tu turno todavía")
    }

    mano := &mesa.ManosJugadores[indiceMano]

    if mano.Estado != string(models.EnCurso) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("esta mano ya fue resuelta")
    }

    if mano.CartasSplit != "" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("ya has hecho split anteriormente")
    }

    cartas := StringToCartas(mano.Cartas)
    if len(cartas) != 2 {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("necesitas exactamente dos cartas para dividir")
    }

    if CalcularValor([]string{cartas[0]}) != CalcularValor([]string{cartas[1]}) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("solo puedes dividir cartas del mismo valor")
    }

    mazo := StringToCartas(mesa.Mazo)
    if len(mazo) < 2 {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no hay suficientes cartas en el mazo para split")
    }

    cartaExtra1, mazo := TomarCarta(mazo)
    cartaExtra2, mazo := TomarCarta(mazo)

    mano1 := []string{cartas[0], cartaExtra1}
    mano2 := []string{cartas[1], cartaExtra2}

    mano.Cartas = CartasToString(mano1)
    mano.CartasSplit = CartasToString(mano2)
    mano.ApuestaSplit = mano.Apuesta
    mano.ManoActual = 1
    mesa.Mazo = CartasToString(mazo)

    valor1 := CalcularValor(mano1)
    if valor1 > 21 {
        mano.ResultadoMano1 = string(models.Perdida)
    }

    s.avanzarManoOMesa(mano, mesa)

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no se pudo dividir la mano: %w", err)
    }

    return s.EstadoParaJugador(mesa, userID), nil
}

func (s *BlackjackService) Seguro(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    if mesa.Estado != "en_curso" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mesa ya ha finalizado")
    }

    indiceMano, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if mesa.JugadorActual != indiceMano {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es tu turno todavía")
    }

    mano := &mesa.ManosJugadores[indiceMano]

    if mano.Estado != string(models.EnCurso) {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("la mano ya no está en curso")
    }

    cartas := StringToCartas(mano.Cartas)
    if len(cartas) != 2 || mano.CartasSplit != "" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("el seguro solo puede tomarse con dos cartas y sin haber hecho split")
    }

    cartasBanca := StringToCartas(mesa.CartasBanca)
    if len(cartasBanca) == 0 || strings.ToUpper(strings.TrimSpace(cartasBanca[0])) != "A" {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no es posible tomar seguro: la banca no muestra un As")
    }

    apuestaSeguro := mano.Apuesta / 2
    mano.Seguro = apuestaSeguro

    if TieneBlackjack(cartasBanca) {
        mano.Seguro = apuestaSeguro * 2
        mano.ResultadoMano1 = string(models.Perdida)
        mano.Estado = CalcularEstadoFinal(mano.ResultadoMano1, "")
        if err := finalizarMesa(mesa); err != nil {
            return dto.BlackjackEstadoDTO{}, fmt.Errorf("fallo al finalizar la mesa: %w", err)
        }
    } else {
        mano.Seguro = 0
    }

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no se pudo actualizar la mesa con el seguro: %w", err)
    }

    return s.EstadoParaJugador(mesa, userID), nil
}

