package ruleta

type NumeroRuleta struct {
	Valor int
	Color string
}

var tableroRuleta = [filasEnTablero][columnasEnTablero]NumeroRuleta{
	{{1, "rojo"}, {2, "negro"}, {3, "rojo"}},
	{{4, "negro"}, {5, "rojo"}, {6, "negro"}},
	{{7, "rojo"}, {8, "negro"}, {9, "rojo"}},
	{{10, "negro"}, {11, "negro"}, {12, "rojo"}},
	{{13, "negro"}, {14, "rojo"}, {15, "negro"}},
	{{16, "rojo"}, {17, "negro"}, {18, "rojo"}},
	{{19, "rojo"}, {20, "negro"}, {21, "rojo"}},
	{{22, "negro"}, {23, "rojo"}, {24, "negro"}},
	{{25, "rojo"}, {26, "negro"}, {27, "rojo"}},
	{{28, "negro"}, {29, "negro"}, {30, "rojo"}},
	{{31, "negro"}, {32, "rojo"}, {33, "negro"}},
	{{34, "rojo"}, {35, "negro"}, {36, "rojo"}},
}

var numeroCero = NumeroRuleta{Valor: 0, Color: "verde"}
