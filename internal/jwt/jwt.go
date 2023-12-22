package jwt

import (
	"homework_platform/internal/bootstrap"
	"strings"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type MyCustomClaims struct {
	ID uint `json:"ID"`
	jwt.RegisteredClaims
}

func CreateToken(id uint) (string, error) {
	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), MyCustomClaims{
		id,
		jwt.RegisteredClaims{
			// ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
	}).SignedString([]byte(bootstrap.Config.JWTSigningString))
}

func GetTokenStr(c *gin.Context) string {
	tokenStr := ""
	// // log.Println(c.Request.Header.Get("Authorization"))
	// // log.Println(c.Request.URL)
	if c.Request.URL.Path == "/api/v1/ai/spark" && c.Request.Method == "GET" {
		tokenStr = c.Query("token")
	} else {
		tokenStr = strings.ReplaceAll(c.Request.Header.Get("Authorization"), "Bearer ", "")
	}
	return tokenStr
}

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

func DecodeTokenStr(tokenStr string) (*jwt.Token, error) {
	var token *jwt.Token
	var err error
	// // log.Println("Decoding", tokenStr)
	at(time.Unix(0, 0), func() {
		token, err = jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(bootstrap.Config.JWTSigningString), nil
		})
	})
	if err != nil {
		return token, err
	}
	return token, nil
}

func MustGetClaims(c *gin.Context) *MyCustomClaims {
	// log.Println("[MustGetClaims]")
	tokenStr := GetTokenStr(c)
	// log.Printf("[MustGetClaims] tokenStr: %s\n", tokenStr)
	token, _ := DecodeTokenStr(tokenStr)
	return token.Claims.(*MyCustomClaims)
}
