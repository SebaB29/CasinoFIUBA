package middleware

import (
	"casino/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(requiredRole ...string) gin.HandlerFunc {
	role := "" // valor por defecto: no se requiere rol
	if len(requiredRole) > 0 {
		role = requiredRole[0]
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header requerido"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header formato inválido"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// Validar rol si se requiere uno
		if role != "" && claims.Rol != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "se requiere rol " + role})
			c.Abort()
			return
		}

		// Podés guardar los claims en el contexto para usar en handlers
		c.Set("userID", claims.UserID)
		c.Set("rol", claims.Rol)

		c.Next()
	}
}
