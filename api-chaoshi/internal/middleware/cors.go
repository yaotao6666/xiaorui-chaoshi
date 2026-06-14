package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedLocalOrigins = map[string]struct{}{
	"http://localhost:3000": {},
	"http://127.0.0.1:3000": {},
	"http://localhost:5173": {},
	"http://127.0.0.1:5173": {},
	"http://localhost:5174": {},
	"http://127.0.0.1:5174": {},
}

// CORS 为本地 H5 和后台调试提供跨域响应头。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && isAllowedOrigin(origin) {
			headers := c.Writer.Header()
			headers.Set("Access-Control-Allow-Origin", origin)
			headers.Set("Access-Control-Allow-Credentials", "true")
			headers.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept, Origin, Cache-Control, X-Requested-With")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			headers.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
			headers.Set("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string) bool {
	if _, ok := allowedLocalOrigins[origin]; ok {
		return true
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := strings.ToLower(parsed.Hostname())
	return host == "localhost" || host == "127.0.0.1"
}
