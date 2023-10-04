package middlewares

import (
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

// 自定义中间件，用于禁用日志记录
func DisableLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将 gin 的 Logger 替换为 ioutil.Discard，即将日志输出到空设备，不记录
		gin.DefaultWriter = ioutil.Discard
		c.Next()
		// 恢复 gin 的 Logger
		gin.DefaultWriter = os.Stdout
	}
}