package errores

import "errors"

var (
	ErrSaldoInsuficiente = errors.New("no tienes suficiente saldo para realizar esta acción")
)
