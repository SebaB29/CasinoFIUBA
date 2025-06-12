package buscaminas

import (
    "math/rand"
    "time"
)

func (t *Tablero) GenerarMinas(cantidad int) {
    rand.Seed(time.Now().UnixNano())

    total := t.Filas * t.Columnas
    if cantidad > total {
        cantidad = total
    }

    posiciones := make(map[int]bool)

    for len(posiciones) < cantidad {
        index := rand.Intn(total)
        if !posiciones[index] {
            posiciones[index] = true
            celda := &t.Celdas[index]
            celda.TieneMina = true
        }
    }
}
