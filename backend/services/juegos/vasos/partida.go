package vasos

import (
	"math/rand"
	"time"
)

type EstadoPartida string

const (
	EnCurso EstadoPartida = "en_curso"
	Ganada  EstadoPartida = "ganada"
	Perdida EstadoPartida = "perdida"
)

type PartidaVasos struct {
	ID               uint          
	PosicionCorrecta int   // 0, 1 o 2         
	Apuesta          float64
	Estado           EstadoPartida
	CreadoEn         time.Time
}

func NuevaPartidaVasos(apuesta float64) *PartidaVasos {
	posicion := rand.Intn(3) // 0, 1 o 2
	return &PartidaVasos{
		PosicionCorrecta: posicion,
		Apuesta:          apuesta,
		Estado:           EnCurso,
		CreadoEn:         time.Now(),
	}
}

func (p *PartidaVasos) Jugar(posicionElegida int) bool {
	if p.Estado != EnCurso {
		return false
	}

	if posicionElegida == p.PosicionCorrecta {
		p.Estado = Ganada
		return true
	}
	p.Estado = Perdida
	return false
}
