package errores

import "errors"

var (
	ErrSaldoInsuficiente = errors.New("no tienes saldo suficiente para realizar esta acción")
	ErrMontoInsuficiente = errors.New("el monto es insuficiente, el monto minimo es 1")
)
