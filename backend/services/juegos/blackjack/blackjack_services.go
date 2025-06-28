package blackjack

import (
    dto "casino/dto/juegos"
    "casino/models"
    repo "casino/repositories/juegos"
    "fmt"
)

// BlackjackService coordina mesas y el hub para eventos
type BlackjackService struct {
    enJuego *BlackjackEnJuego
    Hub     *BlackjackHub
}

// NuevoBlackjackService crea el servicio con su mapa en memoria y su hub
func NuevoBlackjackService() *BlackjackService {
    enJuego := NuevoBlackjackEnJuego()
    hub := NuevoBlackjackHub(enJuego)
    go hub.Run()

    return &BlackjackService{
        enJuego: enJuego,
        Hub:     hub,
    }
}

// CrearMesaInicial arma una mesa con una mano inicial para el creador
func CrearMesaInicial(userID uint, apuesta float64) (*models.MesaBlackjack, error) {
    if apuesta <= 0 {
        return nil, fmt.Errorf("la apuesta debe ser mayor a cero")
    }

    mazo := MezclarMazo(CantidadMazos)

    // Reparto inicial
    carta1Jugador, mazo := TomarCarta(mazo)
    carta2Jugador, mazo := TomarCarta(mazo)
    carta1Banca, mazo := TomarCarta(mazo)
    carta2Banca, mazo := TomarCarta(mazo)

    cartasJugador := []string{carta1Jugador, carta2Jugador}
    cartasBanca := []string{carta1Banca, carta2Banca}

    mano := models.ManoJugadorBlackjack{
        UserID:     userID,
        Apuesta:    apuesta,
        Cartas:     CartasToString(cartasJugador),
        ManoActual: 1,
        Estado:     string(models.EnCurso),
    }

    mesa := &models.MesaBlackjack{
        CartasBanca:    CartasToString(cartasBanca),
        Mazo:           CartasToString(mazo),
        Estado:         "en_curso",
        JugadorActual:  0,
        ManosJugadores: []models.ManoJugadorBlackjack{mano},
    }

    return mesa, nil
}

// AgregarManoJugador agrega un jugador a una mesa existente
func (s *BlackjackService) AgregarManoJugador(mesa *models.MesaBlackjack, userID uint, apuesta float64) error {
    if mesa.Estado != "en_curso" {
        return fmt.Errorf("no se pueden agregar jugadores a una mesa finalizada")
    }

    if len(mesa.ManosJugadores) >= 6 {
        return fmt.Errorf("la mesa ya tiene el máximo de 6 jugadores")
    }

    for _, mano := range mesa.ManosJugadores {
        if mano.UserID == userID {
            return fmt.Errorf("el jugador ya está en la mesa")
        }
    }

    mazo := StringToCartas(mesa.Mazo)
    if len(mazo) < 2 {
        return fmt.Errorf("no hay suficientes cartas para repartir")
    }

    carta1, mazo := TomarCarta(mazo)
    carta2, mazo := TomarCarta(mazo)
    cartas := []string{carta1, carta2}

    nuevaMano := models.ManoJugadorBlackjack{
        UserID:     userID,
        Apuesta:    apuesta,
        Cartas:     CartasToString(cartas),
        ManoActual: 1,
        Estado:     string(models.EnCurso),
    }

    mesa.ManosJugadores = append(mesa.ManosJugadores, nuevaMano)
    mesa.Mazo = CartasToString(mazo)

    return nil
}

// NuevaMesa crea una nueva mesa y registra la primera mano del creador
func (s *BlackjackService) NuevaMesa(userID uint, apuesta float64) (dto.BlackjackEstadoDTO, error) {
    mesa, err := CrearMesaInicial(userID, apuesta)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if err := repo.CrearMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no se pudo persistir la mesa: %w", err)
    }

    s.enJuego.Mutex.Lock()
    s.enJuego.Mesas[mesa.ID] = mesa
    s.enJuego.Mutex.Unlock()

    estado := s.EstadoParaJugador(mesa, userID)
    s.Hub.BroadcastEstado(userID, estado)

    return estado, nil
}

// UnirseAMesa une a un jugador a una mesa ya creada
func (s *BlackjackService) UnirseAMesa(userID uint, idMesa uint, apuesta float64) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada: %w", err)
    }

    if err := s.AgregarManoJugador(mesa, userID, apuesta); err != nil {
        return dto.BlackjackEstadoDTO{}, err
    }

    if err := repo.ActualizarMesa(mesa); err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("no se pudo actualizar la mesa: %w", err)
    }

    estado := s.EstadoParaJugador(mesa, userID)
    s.Hub.BroadcastEstado(userID, estado)

    return estado, nil
}

// ObtenerMesa devuelve una mesa específica con sus manos
func (s *BlackjackService) ObtenerMesa(idMesa uint) (*models.MesaBlackjack, error) {
	return repo.ObtenerMesaConManos(idMesa)
}

// ObtenerEstadoMesa devuelve el estado completo para un jugador
func (s *BlackjackService) ObtenerEstadoMesa(idMesa uint, userID uint) (dto.BlackjackEstadoDTO, error) {
    mesa, err := repo.ObtenerMesaConManos(idMesa)
    if err != nil {
        return dto.BlackjackEstadoDTO{}, fmt.Errorf("mesa no encontrada")
    }

    return s.EstadoParaJugador(mesa, userID), nil
}

// AvanzarTurno avanza al siguiente jugador activo o finaliza la mesa si ya no quedan jugadores en curso
func (s *BlackjackService) AvanzarTurno(mesa *models.MesaBlackjack) {
    total := len(mesa.ManosJugadores)
    for i := mesa.JugadorActual + 1; i < total; i++ {
        if mesa.ManosJugadores[i].Estado == string(models.EnCurso) {
            mesa.JugadorActual = i
            return
        }
    }
    // Si no hay más jugadores con manos activas
    mesa.JugadorActual = total
}

// IndiceManoUsuario devuelve el índice de la mano del jugador en la mesa
func (s *BlackjackService) IndiceManoUsuario(mesa *models.MesaBlackjack, userID uint) (int, error) {
    for i, mano := range mesa.ManosJugadores {
        if mano.UserID == userID {
            return i, nil
        }
    }
    return -1, fmt.Errorf("no se encontró mano para el jugador %d", userID)
}

// EstadoParaJugador prepara el estado visible para el jugador específico
func (s *BlackjackService) EstadoParaJugador(mesa *models.MesaBlackjack, userID uint) dto.BlackjackEstadoDTO {
    indice, err := s.IndiceManoUsuario(mesa, userID)
    if err != nil {
        return dto.BlackjackEstadoDTO{
            Estado: "error: jugador no pertenece a la mesa",
        }
    }

    mano := mesa.ManosJugadores[indice]
    cartasBanca := StringToCartas(mesa.CartasBanca)

    var cartasBancaVisible []string
    var valorBancaVisible, valorBancaTotal int

    if mesa.Estado == "finalizada" {
        cartasBancaVisible = cartasBanca
        valorBancaVisible = CalcularValor(cartasBanca)
        valorBancaTotal = valorBancaVisible
    } else {
        cartasBancaVisible = []string{cartasBanca[0], "???"}
        valorBancaVisible = CalcularValor([]string{cartasBanca[0]})
        valorBancaTotal = 0
    }

    return dto.BlackjackEstadoDTO{
        IDMesa:            mesa.ID,
        JugadorActual:     mesa.JugadorActual,
        EsTurno:           indice == mesa.JugadorActual,
        CartasMano1:       StringToCartas(mano.Cartas),
        ValorMano1:        CalcularValor(StringToCartas(mano.Cartas)),
        CartasMano2:       StringToCartas(mano.CartasSplit),
        ValorMano2:        CalcularValor(StringToCartas(mano.CartasSplit)),
        CartasBanca:       cartasBancaVisible,
        ValorBancaVisible: valorBancaVisible,
        ValorBancaTotal:   valorBancaTotal,
        Estado:            mano.Estado,
        ResultadoMano1:    mano.ResultadoMano1,
        ResultadoMano2:    mano.ResultadoMano2,
        Rendida:           mano.Estado == string(models.Rendida),
        Doblada:           mano.ApuestaSplit == 0 && len(StringToCartas(mano.Cartas)) == 3,
        Seguro:            mano.Seguro,
        Mensaje:           "", 
    }
}

