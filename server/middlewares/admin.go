package middlewares

import (
	"homework_platform/internal/jwt"
	"homework_platform/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		// log.Println("[middleware/AdminCheck]")
		claims := jwt.MustGetClaims(c)
		// log.Println(claims)

		if user, err := models.GetUserByID(claims.ID); err == nil && user.IsAdmin {
			// c.Set("isAdmin", true)
			c.Next()
			return
		}

		c.Status(http.StatusForbidden)
		c.Abort()
	}
}
