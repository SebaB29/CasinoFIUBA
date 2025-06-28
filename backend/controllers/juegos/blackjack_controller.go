package controllers

import (
    "net/http"
    "strconv"
    dto "casino/dto/juegos"
    "casino/services/juegos/blackjack"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    juegos "casino/websocket/juegos"
)

type BlackjackController struct {
    Service *blackjack.BlackjackService
}

func NuevoBlackjackController() *BlackjackController {
    return &BlackjackController{
        Service: blackjack.NuevoBlackjackService(),
    }
}

// NuevaMesa: Crea una nueva mesa para el jugador que envía la request
func (ctrl *BlackjackController) NuevaMesa(c *gin.Context) {
    var input dto.IniciarBlackjackDTO
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetUint("userID")
    res, err := ctrl.Service.NuevaMesa(userID, input.Apuesta)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// UnirseAMesa: Permite que un jugador se una a una mesa existente
func (ctrl *BlackjackController) UnirseAMesa(c *gin.Context) {
    var input dto.UnirseAMesaDTO
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetUint("userID")
    res, err := ctrl.Service.UnirseAMesa(userID, input.IDMesa, input.Apuesta)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// EstadoMesa: Devuelve el estado global de la mesa para el jugador que consulta
func (ctrl *BlackjackController) EstadoMesa(c *gin.Context) {
    idMesa, err := strconv.Atoi(c.Param("id_mesa"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mesa inválido"})
        return
    }

    userID := c.GetUint("userID")
    res, err := ctrl.Service.ObtenerEstadoMesa(uint(idMesa), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// Handler es un helper para acciones como hit, stand, etc.
func (ctrl *BlackjackController) Handler(action func(uint, uint) (gin.H, error)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input dto.JugadaBlackjackDTO
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        userID := c.GetUint("userID")
        res, err := action(input.IDMesa, userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, res)
    }
}

// JugarWS: Handler que acepta la conexión WebSocket
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
