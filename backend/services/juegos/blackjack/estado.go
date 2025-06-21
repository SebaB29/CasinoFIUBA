package blackjack

import "casino/models"

func EstadoFromResultado(resultado string) models.EstadoBlackjack {
	switch resultado {
	case "ganada":
		return models.Ganada
	case "perdida":
		return models.Perdida
	case "empatada":
		return models.Empatada
	default:
		return models.EnCurso
	}
}
