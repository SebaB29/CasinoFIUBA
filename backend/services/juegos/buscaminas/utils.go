package buscaminas

func CoordenadaValida(x, y, maxX, maxY int) bool {
    return x >= 0 && x < maxX && y >= 0 && y < maxY
}
