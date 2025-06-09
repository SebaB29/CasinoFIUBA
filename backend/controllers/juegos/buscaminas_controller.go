package controllers

import (
	"bytes"
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
	Fila      *int `json:"fila" binding:"required,min=0,max=4"`
	Col       *int `json:"col" binding:"required,min=0,max=4"`
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
	log.Printf("Intentando crear nueva partida - UserID: %d, Minas: %d, Apuesta: %.2f", userID, req.Minas, req.Apuesta)

	partida, err := buscaminas.CrearPartida(5, 5, req.Minas, req.Apuesta)
	if err != nil {
		log.Println("Error al crear partida:", err)
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
		log.Println("Error al guardar partida en base de datos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la partida"})
		return
	}

	log.Printf("Partida creada correctamente (ID: %d) para el usuario %d", partidaDB.ID, userID)

	c.JSON(http.StatusCreated, gin.H{
		"id_partida":      partidaDB.ID,
		"estado":          partida.Estado,
		"minas":           partida.CantidadMinas,
		"monto_apostado":  partida.MontoApostado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})
}

func AbrirCeldaBuscaminas(c *gin.Context) {
	// Loguear el body crudo
	rawData, _ := io.ReadAll(c.Request.Body)
	log.Println("🟡 Datos crudos recibidos en /buscaminas/abrir:", string(rawData))
	// Resetear el Body para que se pueda usar otra vez
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	// Bindeo normal
	var req AbrirCeldaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("🔴 Error de request al abrir celda:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("🟢 Usuario %d intenta abrir celda (%d, %d) en partida %d", userID, req.Fila, req.Col, req.IDPartida)

	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, req.IDPartida).Error; err != nil {
		log.Println("🔴 Partida no encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	partida, err := buscaminas.CrearPartida(5, 5, partidaDB.CantidadMinas, partidaDB.MontoApostado)
	if err != nil {
		log.Println("🔴 Error reconstruyendo partida:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}
	partida.CeldasAbiertas = partidaDB.CeldasAbiertas

	err = partida.AbrirCelda(*req.Fila, *req.Col)
	if err != nil {
		log.Println("🟠 Resultado: celda con mina o inválida:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partidaDB.Estado = string(partida.Estado)
	partidaDB.CeldasAbiertas = partida.CeldasAbiertas
	db.DB.Save(&partidaDB)

	log.Printf("✅ Celda abierta con éxito - Celdas abiertas: %d | Estado: %s", partida.CeldasAbiertas, partida.Estado)

	c.JSON(http.StatusOK, gin.H{
		"estado":          partida.Estado,
		"celdas_abiertas": partida.CeldasAbiertas,
		"tablero":         ocultarMinas(partida.Tablero),
	})
}

func RetirarseBuscaminas(c *gin.Context) {
	var req RetirarseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error de request al retirarse:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	log.Printf("Usuario %d se retira de la partida %d", userID, req.IDPartida)

	var partidaDB models.PartidaBuscaminas
	if err := db.DB.First(&partidaDB, req.IDPartida).Error; err != nil {
		log.Println("Partida no encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Partida no encontrada"})
		return
	}

	partida, err := buscaminas.CrearPartida(5, 5, partidaDB.CantidadMinas, partidaDB.MontoApostado)
	if err != nil {
		log.Println("Error reconstruyendo partida:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reconstruyendo la partida"})
		return
	}
	partida.CeldasAbiertas = partidaDB.CeldasAbiertas

	premio, err := partida.Retirarse()
	if err != nil {
		log.Println("No se pudo retirar:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partidaDB.Estado = string(partida.Estado)
	db.DB.Save(&partidaDB)

	log.Printf("Usuario %d se retiró con premio: %.2f", userID, premio)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Te retiraste exitosamente",
		"estado":  partida.Estado,
		"premio":  premio,
	})
}

func ocultarMinas(t *buscaminas.Tablero) [][]map[string]interface{} {
	resultado := make([][]map[string]interface{}, t.Filas)
	for i := range t.Celdas {
		resultado[i] = make([]map[string]interface{}, t.Columnas)
		for j, celda := range t.Celdas[i] {
			resultado[i][j] = map[string]interface{}{
				"abierta":        celda.Abierta,
				"marcada":        celda.Marcada,
				"minas_cercanas": celda.MinasVecinas,
			}
		}
	}
	return resultado
}
