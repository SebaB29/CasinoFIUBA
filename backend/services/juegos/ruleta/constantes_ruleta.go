package ruleta

// Tablero
const (
	filasEnTablero    = 12
	columnasEnTablero = 3
)

// Rango de números en la ruleta
const (
	NumeroMinimoRuleta    = 0
	NumeroMaximoRuleta    = 36
	CantidadNumerosRuleta = 37
)

// Rango para docenas
const (
	PrimeraDocena = 1
	SegundaDocena = 2
	TerceraDocena = 3

	MinDocena1 = 1
	MaxDocena1 = 12
	MinDocena2 = 13
	MaxDocena2 = 24
	MinDocena3 = 25
	MaxDocena3 = 36
)

// Rango para apuestas "Alto/Bajo"
const (
	MinBajo = 1
	MaxBajo = 18
	MinAlto = 19
	MaxAlto = 36
)

// Multiplicadores de pago según tipo de apuesta
const (
	MultiplicadorPleno    = 36.0
	MultiplicadorDividida = 18.0
	MultiplicadorCalle    = 12.0
	MultiplicadorCuadro   = 9.0
	MultiplicadorDocena   = 3.0
	MultiplicadorSimple   = 2.0
	SinGanancia           = 0.0
)

// Cantidad de números que abarca cada tipo de apuesta
const (
	CantidadNumerosPleno    = 1
	CantidadNumerosDividida = 2
	CantidadNumerosCalle    = 3
	CantidadNumerosCuadro   = 4
)

// Valores válidos para apuestas de color, paridad y alto/bajo
var (
	ColoresValidos   = map[string]bool{"rojo": true, "negro": true}
	ParidadesValidas = map[string]bool{"par": true, "impar": true}
	AltoBajoValidos  = map[string]bool{"alto": true, "bajo": true}
)