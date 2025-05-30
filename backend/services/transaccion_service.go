package services

import (
	"casino/db"
	"casino/errores"
	"casino/models"
	"casino/repositories"
)

const MontoMinimo = 1.0

type TransaccionServiceInterface interface {
	Depositar(usuarioID uint, monto float64) error
	Extraer(usuarioID uint, monto float64) error
}

type TransaccionService struct {
	usuarioRepository     repositories.UsuarioRepositoryInterface
	transaccionRepository repositories.TransaccionRepositoryInterface
}

func NewTransaccionService() *TransaccionService {
	return &TransaccionService{
		usuarioRepository:     repositories.NewUsuarioRepository(db.DB),
		transaccionRepository: repositories.NewTransaccionRepository(db.DB),
	}
}

func (transaccionService *TransaccionService) Depositar(usuarioID uint, monto float64) error {
	usuario, err := transaccionService.usuarioRepository.ObtenerPorID(usuarioID)
	if err != nil || usuario == nil {
		return errores.ErrUsuarioNoEncontrado
	}

	if monto < MontoMinimo {
		return errores.ErrMontoInsuficiente
	}

	usuario.Saldo += monto
	if err := transaccionService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	transaccion := &models.Transaccion{
		UsuarioID: usuario.ID,
		Tipo:      "deposito",
		Monto:     monto,
	}

	return transaccionService.transaccionRepository.Crear(transaccion)
}

func (transaccionService *TransaccionService) Extraer(usuarioID uint, monto float64) error {
	usuario, err := transaccionService.usuarioRepository.ObtenerPorID(usuarioID)
	if err != nil || usuario == nil {
		return errores.ErrUsuarioNoEncontrado
	}

	if monto < MontoMinimo {
		return errores.ErrMontoInsuficiente
	}

	if usuario.Saldo < monto {
		return errores.ErrSaldoInsuficiente
	}

	usuario.Saldo -= monto
	if err := transaccionService.usuarioRepository.Actualizar(usuario); err != nil {
		return err
	}

	transaccion := &models.Transaccion{
		UsuarioID: usuario.ID,
		Tipo:      "extraccion",
		Monto:     monto,
	}

	return transaccionService.transaccionRepository.Crear(transaccion)
}
