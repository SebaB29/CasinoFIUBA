package buscaminas

// CoordenadaValida verifica si una posición (fila, columna) está dentro de los límites del tablero.
func CoordenadaValida(fila, columna, maxFilas, maxColumnas int) bool {
	return fila >= 0 && fila < maxFilas && columna >= 0 && columna < maxColumnas
}

// IndexToCoords convierte un índice lineal a coordenadas (fila, columna) dentro del tablero.
func IndexToCoords(index, columnas int) (int, int) {
	return index / columnas, index % columnas
}

// Vecinos devuelve una lista de coordenadas vecinas válidas para una celda dada.
func Vecinos(fila, columna, maxFilas, maxColumnas int) [][2]int {
	vecinos := make([][2]int, 0)

	for df := -1; df <= 1; df++ {
		for dc := -1; dc <= 1; dc++ {
			if df == 0 && dc == 0 {
				continue // ignorar la celda actual
			}
			nf, nc := fila+df, columna+dc
			if CoordenadaValida(nf, nc, maxFilas, maxColumnas) {
				vecinos = append(vecinos, [2]int{nf, nc})
			}
		}
	}

	return vecinos
}
