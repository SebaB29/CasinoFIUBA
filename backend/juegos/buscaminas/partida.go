package buscaminas

import (
	"errors"
)

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
	MinasRestantes  int
	CeldasAbiertas  int
	CantidadMinas   int
	MontoApostado   float64
}

// CrearPartida inicializa una partida nueva delegando la colocación de minas
func CrearPartida(filas, columnas, cantidadMinas int, montoApostado float64) (*Partida, error) {
	if cantidadMinas >= filas*columnas {
		return nil, errors.New("la cantidad de minas excede el tamaño del tablero")
	}

	tablero := NewTablero(filas, columnas)
	tablero.GenerarMinas(cantidadMinas)
	tablero.ContarMinasVecinas()

	partida := &Partida{
		Tablero:         tablero,
		Estado:          EnCurso,
		MinasRestantes:  cantidadMinas,
		CantidadMinas:   cantidadMinas,
		CeldasAbiertas:  0,
		MontoApostado:   montoApostado,
	}

	return partida, nil
}

func (p *Partida) AbrirCelda(f, c int) error {
	if p.Estado != EnCurso {
		return errors.New("la partida ya ha finalizado")
	}

	if f < 0 || f >= p.Tablero.Filas || c < 0 || c >= p.Tablero.Columnas {
		return errors.New("posición inválida")
	}

	celda := p.Tablero.Celdas[f][c]
	if celda.Abierta {
		return nil
	}

	if celda.TieneMina {
		p.Estado = Perdida
		p.Tablero.Celdas[f][c].Abierta = true
		return errors.New("¡boom! pisaste una mina")
	}

	p.abrirRecursivo(f, c)

	totalCeldas := p.Tablero.Filas * p.Tablero.Columnas
	if p.CeldasAbiertas == totalCeldas-p.CantidadMinas {
		p.Estado = Ganada
	}

	return nil
}

func (p *Partida) abrirRecursivo(f, c int) {
	celda := p.Tablero.Celdas[f][c]
	if celda.Abierta || celda.TieneMina {
		return
	}

	p.Tablero.Celdas[f][c].Abierta = true
	p.CeldasAbiertas++

	if celda.MinasVecinas == 0 {
		for df := -1; df <= 1; df++ {
			for dc := -1; dc <= 1; dc++ {
				if df != 0 || dc != 0 {
					nf, nc := f+df, c+dc
					if CoordenadaValida(nf, nc, p.Tablero.Filas, p.Tablero.Columnas) {
						p.abrirRecursivo(nf, nc)
					}
				}
			}
		}
	}
}

func (p *Partida) Retirarse() (float64, error) {
	if p.Estado != EnCurso {
		return 0, errors.New("la partida ya ha finalizado")
	}

	p.Estado = Retirada

	totalCeldas := p.Tablero.Filas * p.Tablero.Columnas
	celdasJugables := totalCeldas - p.CantidadMinas
	if celdasJugables == 0 {
		return 0, errors.New("no hay celdas jugables")
	}

	progreso := float64(p.CeldasAbiertas) / float64(celdasJugables)
	factorDificultad := 0.2 * float64(p.CantidadMinas)

	premio := p.MontoApostado * (1 + factorDificultad*progreso)
	return premio, nil
}
