package dto

type IniciarBlackjackDTO struct {
	Apuesta float64 `json:"apuesta" binding:"required"`
}

type UnirseAMesaDTO struct {
	IDMesa  uint    `json:"id_mesa" binding:"required"`
	Apuesta float64 `json:"apuesta" binding:"required"`
}

type JugadaBlackjackDTO struct {
	IDMesa uint `json:"id_mesa" binding:"required"`
}

type SplitBlackjackDTO struct {
	IDMesa uint `json:"id_mesa" binding:"required"`
}

type SeguroBlackjackDTO struct {
	IDMesa uint `json:"id_mesa" binding:"required"`
}

type BlackjackEstadoDTO struct {
	IDMesa             uint     `json:"id_mesa"`
	JugadorActual      int      `json:"jugador_actual"`
	EsTurno            bool     `json:"es_turno"`
	CartasMano1        []string `json:"cartas_mano_1"`
	ValorMano1         int      `json:"valor_mano_1"`
	CartasMano2        []string `json:"cartas_mano_2"`
	ValorMano2         int      `json:"valor_mano_2"`
	CartasBanca        []string `json:"cartas_banca"`
	ValorBancaVisible  int      `json:"valor_banca_visible"`
	Estado             string   `json:"estado"` 
	ResultadoMano1     string `json:"resultado_mano_1"`
    ResultadoMano2     string `json:"resultado_mano_2"`
    Rendida            bool   `json:"rendida"`
    Doblada            bool   `json:"doblada"`
    Seguro             float64 `json:"seguro"`
    Mensaje            string `json:"mensaje"`
    ValorBancaTotal    int    `json:"valor_banca_total"`
}
