package services

var ruletaTablero = [12][3]int{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 9},
	{10, 11, 12},
	{13, 14, 15},
	{16, 17, 18},
	{19, 20, 21},
	{22, 23, 24},
	{25, 26, 27},
	{28, 29, 30},
	{31, 32, 33},
	{34, 35, 36},
}

func posicionEnTablero(n int) (fila, col int, ok bool) {
	for i, filaArr := range ruletaTablero {
		for j, v := range filaArr {
			if v == n {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}

func SonAdyacentes(a, b int) bool {
	fa, ca, oka := posicionEnTablero(a)
	fb, cb, okb := posicionEnTablero(b)

	if !oka || !okb {
		return false
	}

	df, dc := fa-fb, ca-cb
	return (df == 0 && abs(dc) == 1) || (abs(df) == 1 && dc == 0)
}

func EsCalleValida(numeros []int) bool {
	for _, fila := range ruletaTablero {
		if contieneTodos(fila[:], numeros) {
			return true
		}
	}
	return false
}

func EsCuadroValido(numeros []int) bool {
	for i := 0; i < len(ruletaTablero)-1; i++ {
		for j := 0; j < 2; j++ {
			cuadro := []int{
				ruletaTablero[i][j], ruletaTablero[i][j+1],
				ruletaTablero[i+1][j], ruletaTablero[i+1][j+1],
			}
			if contieneTodos(cuadro, numeros) {
				return true
			}
		}
	}
	return false
}

func contieneTodos(arr []int, target []int) bool {
	m := make(map[int]bool)
	for _, n := range arr {
		m[n] = true
	}
	for _, n := range target {
		if !m[n] {
			return false
		}
	}
	return true
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
