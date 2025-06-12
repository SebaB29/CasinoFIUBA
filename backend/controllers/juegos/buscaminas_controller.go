package controllers

import (
	"bytes"
	"casino/config"
	"casino/db"
	"casino/juegos/buscaminas"
	"casino/models"
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
		log.Println("Error de request al crear partida:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("🛠️ Intentando crear nueva partida - UserID: %d, Minas: %d, Apuesta: %.2f", userID, req.Minas, req.Apuesta)

	partida, err := buscaminas.CrearPartida(5, 5, req.Minas, req.Apuesta)
	if err != nil {
		log.Println("❌ Error al crear partida:", err)
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
		log.Println("❌ Error al guardar partida en base de datos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la partida"})
		return
	}

	log.Printf("✅ Partida creada correctamente (ID: %d) para el usuario %d", partidaDB.ID, userID)

	c.JSON(http.StatusCreated, gin.H{
		"id_partida":      partidaDB.ID,
		"estado":          partida.Estado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})
}

func AbrirCeldaBuscaminas(c *gin.Context) {
	rawData, _ := io.ReadAll(c.Request.Body)
	log.Println("📦 Datos crudos recibidos en /buscaminas/abrir:", string(rawData))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	var req AbrirCeldaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("🔴 Error de request al abrir celda:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("🟢 Usuario %d intenta abrir celda (%d, %d) en partida %d", userID, *req.X, *req.Y, req.IDPartida)

	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, req.IDPartida).Error; err != nil {
		log.Println("🔴 Partida no encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	partida, err := buscaminas.CrearPartida(5, 5, partidaDB.CantidadMinas, partidaDB.MontoApostado)
	if err != nil {
		log.Println("❌ Error reconstruyendo partida:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}
	partida.CeldasAbiertas = partidaDB.CeldasAbiertas

	err = partida.AbrirCelda(*req.X, *req.Y)
	if err != nil {
		log.Println("🟠 Resultado: celda con mina o inválida:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partidaDB.Estado = string(partida.Estado)
	partidaDB.CeldasAbiertas = partida.CeldasAbiertas
	db.DB.Save(&partidaDB)

	log.Printf("✅ Celda abierta - Celdas abiertas: %d | Estado: %s", partida.CeldasAbiertas, partida.Estado)

	responderPartida(c, partida, &partidaDB, http.StatusOK)
}

func RetirarseBuscaminas(c *gin.Context) {
	var req RetirarseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Error de request al retirarse:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("📤 Usuario %d se retira de la partida %d", userID, req.IDPartida)

	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, req.IDPartida).Error; err != nil {
		log.Println("🔴 Partida no encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	partida, err := buscaminas.CrearPartida(5, 5, partidaDB.CantidadMinas, partidaDB.MontoApostado)
	if err != nil {
		log.Println("❌ Error reconstruyendo partida:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}
	partida.CeldasAbiertas = partidaDB.CeldasAbiertas

	premio, err := partida.Retirarse()
	if err != nil {
		log.Println("⚠️ No se pudo retirar:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partidaDB.Estado = string(partida.Estado)
	db.DB.Save(&partidaDB)

	log.Printf("💰 Usuario %d se retiró con premio: %.2f", userID, premio)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Te retiraste exitosamente",
		"estado":  partida.Estado,
		"premio":  premio,
	})
}

// Devuelve solo las celdas visibles al usuario
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

// Factor común para devolver estado de partida + tablero oculto
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
