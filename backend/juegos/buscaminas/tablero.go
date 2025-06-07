package buscaminas

type Tablero struct {
	Filas    int
	Columnas int
	Celdas   [][]Celda
}

// ContarMinasVecinas actualiza cada celda con la cantidad de minas alrededor.
func (t *Tablero) ContarMinasVecinas() {
	direcciones := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1},           {0, 1},
		{1, -1},  {1, 0},  {1, 1},
	}

	for fila := 0; fila < t.Filas; fila++ {
		for col := 0; col < t.Columnas; col++ {
			celda := t.Celdas[fila][col]
			if celda.TieneMina {
				continue
			}

			conteo := 0
			for _, dir := range direcciones {
				nf := fila + dir[0]
				nc := col + dir[1]

				if nf >= 0 && nf < t.Filas && nc >= 0 && nc < t.Columnas &&
					t.Celdas[nf][nc].TieneMina {
					conteo++
				}
			}
			t.Celdas[fila][col].MinasVecinas = conteo
		}
	}
}

func NewTablero(filas, columnas int) *Tablero {
	celdas := make([][]Celda, filas)
	for i := 0; i < filas; i++ {
		celdas[i] = make([]Celda, columnas)
	}
	return &Tablero{
		Filas:    filas,
		Columnas: columnas,
		Celdas:   celdas,
	}
}

