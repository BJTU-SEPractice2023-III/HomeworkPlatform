package middlewares

import (
	"homework_platform/internal/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenStr, err := c.Cookie("token")
		// log.Printf("[middlewares/JWTAuth]: Token from cookies: %v\n", tokenStr)
		if err != nil {
			tokenStr = jwt.GetTokenStr(c)
			// log.Printf("[middlewares/JWTAuth]: Token from headers: %v\n", tokenStr)
		}

		token, err := jwt.DecodeTokenStr(tokenStr)
		// log.Println(token, err)

		if err != nil || !token.Valid {
			log.Printf("[middlewares/JWTAuth]: Token not valid: %v\n", err)

			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		id := token.Claims.(*jwt.MyCustomClaims).ID
		c.Set("ID", id)
		c.Next()
	}
}
