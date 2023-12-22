package middlewares

import (
	"homework_platform/internal/bootstrap"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"net/http/httputil"
	"net/url"
)

func Frontend() gin.HandlerFunc {
	ignoreFunc := func(c *gin.Context) {
		c.Next()
	}

	if bootstrap.StaticFS == nil {
		return ignoreFunc
	}

	// 读取index.html
	file, err := bootstrap.StaticFS.Open("/index.html")
	if err != nil {
		// log.Println("Static file \"index.html\" does not exist, it might affect the display of the homepage.")
		return ignoreFunc
	}

	fileContentBytes, err := io.ReadAll(file)
	if err != nil {
		// log.Println("Cannot read static file \"index.html\", it might affect the display of the homepage.")
		return ignoreFunc
	}
	fileContent := string(fileContentBytes)

	fileServer := http.FileServer(bootstrap.StaticFS)
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 跳过
		if strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		if path == "/index.html" || path == "/" || !bootstrap.StaticFS.Exists("/", path) {

			c.Header("Content-Type", "text/html")
			c.String(200, fileContent)
			c.Abort()
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func FrontendReverseProxy() gin.HandlerFunc {
	target, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// API 跳过
		if strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		// 修改请求头等信息
		c.Request.Host = target.Host
		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))
		c.Request.Header.Set("Host", target.Host)

		// 执行反向代理
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
