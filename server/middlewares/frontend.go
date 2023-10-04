package middlewares

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"net/http/httputil"
	"net/url"
)

func Frontend(fs http.FileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 跳过
		if strings.HasPrefix(path, "/api") {
			c.Next()
		} else {
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
	}
}

func FrontendReverseProxy() gin.HandlerFunc {
    target, _ := url.Parse("http://localhost:3000")
    proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		gin.DefaultWriter = io.Discard
		// 修改请求头等信息
		c.Request.Host = target.Host
		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))
		c.Request.Header.Set("Host", target.Host)

		// 执行反向代理
		proxy.ServeHTTP(c.Writer, c.Request)
		gin.DefaultWriter = os.Stdout
	}
}
