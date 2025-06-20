package controllers

import (
	"casino/services/juegos/ruleta"
	ws "casino/websocket/juegos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

type RuletaController struct {
	ruletaService *ruleta.RuletaService
}

func NewRuletaController() *RuletaController {
	return &RuletaController{
		ruletaService: ruleta.NewRuletaService(),
	}
}

func (ctrl *RuletaController) JugarWS(ctx *gin.Context) {
	conexion, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo establecer WebSocket"})
		return
	}

	userID := ctx.GetUint("userID")
	handler := ws.NewRuletaSocketHandler(conexion, userID, ctrl.ruletaService)
	handler.Manejar()
}
