package buscaminas

type Celda struct {
	Abierta   bool
	TieneMina     bool
	Marcada       bool
	MinasVecinas int
}
