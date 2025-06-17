package buscaminas

import "errors"

type EstadoPartida string

const (
    EnCurso  EstadoPartida = "en_curso"
    Ganada   EstadoPartida = "ganada"
    Perdida  EstadoPartida = "perdida"
    Retirada EstadoPartida = "retirada"
)

type Partida struct {
    Tablero         *Tablero
    Estado          EstadoPartida
    CeldasAbiertas  int
    CantidadMinas   int
    MontoApostado   float64
}

func CrearPartida(filas, columnas, cantidadMinas int, montoApostado float64) (*Partida, error) {
    if cantidadMinas >= filas*columnas {
        return nil, errors.New("la cantidad de minas excede el tamaño del tablero")
    }

    tablero := NewTablero(filas, columnas)
    tablero.GenerarMinas(cantidadMinas)

    return &Partida{
        Tablero:        tablero,
        Estado:         EnCurso,
        CantidadMinas:  cantidadMinas,
        MontoApostado:  montoApostado,
        CeldasAbiertas: 0,
    }, nil
}

func (p *Partida) AbrirCelda(x, y int) error {
    if p.Estado != EnCurso {
        return errors.New("la partida ya ha finalizado")
    }

    if !CoordenadaValida(x, y, p.Tablero.Columnas, p.Tablero.Filas) {
        return errors.New("posición inválida")
    }

    celda := p.Tablero.GetCelda(x, y)
    if celda == nil || celda.Abierta {
        return nil
    }

    celda.Abierta = true
    p.CeldasAbiertas++

    if celda.TieneMina {
        p.Estado = Perdida
        return errors.New("¡boom! pisaste una mina")
    }

    totalCeldas := p.Tablero.Filas * p.Tablero.Columnas
    if p.CeldasAbiertas == totalCeldas-p.CantidadMinas {
        p.Estado = Ganada
    }

    return nil
}

func (p *Partida) Retirarse() (float64, error) {
    if p.Estado != EnCurso {
        return 0, errors.New("la partida ya ha finalizado")
    }

    p.Estado = Retirada

    total := p.Tablero.Filas * p.Tablero.Columnas
    jugables := total - p.CantidadMinas
    if jugables == 0 {
        return 0, errors.New("no hay celdas jugables")
    }

    progreso := float64(p.CeldasAbiertas) / float64(jugables)
    dificultad := 0.2 * float64(p.CantidadMinas)

    premio := p.MontoApostado * (1 + dificultad*progreso)
    return premio, nil
}
