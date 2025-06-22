package controllers

import (
	dto "casino/dto/juegos"
	"casino/services/juegos/blackjack"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
	juegos "casino/websocket/juegos"

	"github.com/gin-gonic/gin"
)

type BlackjackController struct {
	Service *blackjack.BlackjackService
}

func NewBlackjackController() *BlackjackController {
	return &BlackjackController{
		Service: blackjack.NewBlackjackService(),
	}
}

func (ctrl *BlackjackController) CrearPartida(c *gin.Context) {
	var input dto.IniciarBlackjackDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("userID")
	res, err := ctrl.Service.CrearPartida(userID, input.Apuesta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *BlackjackController) Handler(action func(uint) (gin.H, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input dto.JugadaBlackjackDTO
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := action(input.IDPartida)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func (ctrl *BlackjackController) Estado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id_partida"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	res, err := ctrl.Service.ObtenerEstado(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *BlackjackController) JugarWS(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo abrir WS"})
		return
	}
	userID := c.GetUint("userID")
	handler := juegos.NewBlackjackSocketHandler(conn, userID, ctrl.Service, ctrl.Service.Hub)
	handler.Manejar()
}