package buscaminas

type Tablero struct {
    Filas    int
    Columnas int
    Celdas   []Celda
}

func NewTablero(filas, columnas int) *Tablero {
    total := filas * columnas
    celdas := make([]Celda, 0, total)
    for y := 0; y < filas; y++ {
        for x := 0; x < columnas; x++ {
            celdas = append(celdas, Celda{X: x, Y: y})
        }
    }
    return &Tablero{
        Filas:    filas,
        Columnas: columnas,
        Celdas:   celdas,
    }
}

func (t *Tablero) GetCelda(x, y int) *Celda {
    for i := range t.Celdas {
        if t.Celdas[i].X == x && t.Celdas[i].Y == y {
            return &t.Celdas[i]
        }
    }
    return nil
}
