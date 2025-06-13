package models

import "time"

type JugadaRuleta struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UsuarioID uint    `json:"usuario_id" gorm:"not null"`
	Usuario   Usuario `json:"-" gorm:"foreignKey:UsuarioID"`

	MontoApostado float64 `json:"monto_apostado"`
	Ganancia      float64 `json:"ganancia"`

	TipoApuesta string `json:"tipo_apuesta" gorm:"type:varchar(20)"`

	Numeros  IntSlice `json:"numeros,omitempty" gorm:"type:text"`
	Color    string   `json:"color,omitempty" gorm:"type:varchar(10)"`
	Paridad  string   `json:"paridad,omitempty" gorm:"type:varchar(5)"`
	AltoBajo string   `json:"alto_bajo,omitempty" gorm:"type:varchar(5)"`
	Docena   int      `json:"docena,omitempty"`

	NumeroGanador int    `json:"numero_ganador"`
	ColorGanador  string `json:"color_ganador"`

	Fecha time.Time `json:"fecha"`
}

func (JugadaRuleta) TableName() string {
	return "jugadas_ruleta"
}
