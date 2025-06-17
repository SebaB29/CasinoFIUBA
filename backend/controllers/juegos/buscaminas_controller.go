package controllers

import (
	"bytes"
	"casino/config"
	"casino/db"
	"casino/models"
	"casino/services/juegos/buscaminas"
	repositories "casino/repositories/juegos"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AbrirCeldaRequest struct {
	IDPartida uint `json:"id_partida" binding:"required"`
	X         *int `json:"x" binding:"required,min=0,max=4"`
	Y         *int `json:"y" binding:"required,min=0,max=4"`
}

type NuevaPartidaRequest struct {
	Minas   int     `json:"minas" binding:"required,min=1,max=24"`
	Apuesta float64 `json:"apuesta" binding:"required,gt=0"`
}

type RetirarseRequest struct {
	IDPartida uint `json:"id_partida" binding:"required"`
}

func CrearPartidaBuscaminas(c *gin.Context) {
	var req NuevaPartidaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Error de request al crear partida:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("🛠️ Creando nueva partida (Usuario %d) - Minas: %d | Apuesta: %.2f", userID, req.Minas, req.Apuesta)

	partida, err := buscaminas.CrearPartida(5, 5, req.Minas, req.Apuesta)
	if err != nil {
		log.Println("❌ Error al crear partida lógica:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repo := repositories.NewBuscaminasRepository(db.DB)
	partidaDB := &models.PartidaBuscaminas{
		UsuarioID:      &userID,
		Estado:         string(partida.Estado),
		CeldasAbiertas: partida.CeldasAbiertas,
		CantidadMinas:  partida.CantidadMinas,
		MontoApostado:  partida.MontoApostado,
	}

	if err := repo.Crear(partidaDB); err != nil {
		log.Println("❌ Error al guardar partida en DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la partida"})
		return
	}

	log.Printf("✅ Partida creada correctamente (ID %d)", partidaDB.ID)

	c.JSON(http.StatusCreated, gin.H{
		"id_partida":      partidaDB.ID,
		"estado":          partida.Estado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})
}

func AbrirCeldaBuscaminas(c *gin.Context) {
	rawData, _ := io.ReadAll(c.Request.Body)
	log.Println("📦 Body recibido en /abrir:", string(rawData))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	var req AbrirCeldaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("🔴 Error al bindear:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("👆 Usuario %d intenta abrir celda (%d,%d) en partida %d", userID, *req.X, *req.Y, req.IDPartida)

	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, req.IDPartida).Error; err != nil {
		log.Println("🔴 Partida no encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	partida, err := reconstruirPartida(&partidaDB)
	if err != nil {
		log.Println("❌ Error reconstruyendo partida:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}

	err = partida.AbrirCelda(*req.X, *req.Y)
	if err != nil {
		// Pisar una mina u otro error del juego es 200 OK
		log.Println("💥 Jugada válida - celda con mina u otra condición:", err)
	}

	partidaDB.Estado = string(partida.Estado)
	partidaDB.CeldasAbiertas = partida.CeldasAbiertas
	db.DB.Save(&partidaDB)

	log.Printf("🟢 Resultado - Estado: %s | Celdas abiertas: %d", partida.Estado, partida.CeldasAbiertas)

	responderPartida(c, partida, &partidaDB, http.StatusOK)
}

func RetirarseBuscaminas(c *gin.Context) {
	var req RetirarseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Error en retiro:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("📤 Usuario %d intenta retirarse de partida %d", userID, req.IDPartida)

	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, req.IDPartida).Error; err != nil {
		log.Println("🔴 Partida no encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	switch partidaDB.Estado {
	case "ganada":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ya ganaste, no podés retirarte"})
		return
	case "perdida":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ya perdiste, no podés retirarte"})
		return
	}

	partida, err := reconstruirPartida(&partidaDB)
	if err != nil {
		log.Println("❌ Error reconstruyendo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}

	premio, err := partida.Retirarse()
	if err != nil {
		log.Println("⚠️ No se pudo retirar:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partidaDB.Estado = string(partida.Estado)
	db.DB.Save(&partidaDB)

	log.Printf("💰 Usuario %d se retiró y ganó: %.2f", userID, premio)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Te retiraste exitosamente",
		"estado":  partida.Estado,
		"premio":  premio,
	})
}

// Helpers

func reconstruirPartida(partidaDB *models.PartidaBuscaminas) (*buscaminas.Partida, error) {
	partida, err := buscaminas.CrearPartida(5, 5, partidaDB.CantidadMinas, partidaDB.MontoApostado)
	if err != nil {
		return nil, err
	}
	partida.CeldasAbiertas = partidaDB.CeldasAbiertas
	partida.Estado = buscaminas.EstadoPartida(partidaDB.Estado)
	return partida, nil
}

func esFinalizada(estado string) bool {
	return estado == "ganada" || estado == "perdida" || estado == "retirada"
}

func ocultarMinas(t *buscaminas.Tablero) []map[string]interface{} {
	resultado := make([]map[string]interface{}, 0, len(t.Celdas))
	for _, celda := range t.Celdas {
		resultado = append(resultado, map[string]interface{}{
			"x":       celda.X,
			"y":       celda.Y,
			"abierta": celda.Abierta,
		})
	}
	return resultado
}

func responderPartida(c *gin.Context, partida *buscaminas.Partida, partidaDB *models.PartidaBuscaminas, status int) {
	c.JSON(status, gin.H{
		"estado":          partida.Estado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})
}

func VerMinasDebug(c *gin.Context) {
	if !config.IsDevMode {
		c.JSON(http.StatusForbidden, gin.H{"error": "Debug deshabilitado en producción"})
		return
	}

	idStr := c.Param("id")
	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, idStr).Error; err != nil {
		log.Println("🔴 Partida no encontrada para debug:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	partida, err := buscaminas.CrearPartida(5, 5, partidaDB.CantidadMinas, partidaDB.MontoApostado)
	if err != nil {
		log.Println("❌ Error reconstruyendo partida para debug:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}
	partida.CeldasAbiertas = partidaDB.CeldasAbiertas

	minas := make([]map[string]int, 0)
	for _, celda := range partida.Tablero.Celdas {
		if celda.TieneMina {
			minas = append(minas, map[string]int{
				"x": celda.X,
				"y": celda.Y,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id_partida": partidaDB.ID,
		"minas":      minas,
	})
}