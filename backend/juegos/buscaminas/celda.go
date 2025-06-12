package buscaminas

type Celda struct {
    X         int  `json:"x"`
    Y         int  `json:"y"`
    Abierta   bool `json:"abierta"`
    TieneMina bool `json:"-"`
}
