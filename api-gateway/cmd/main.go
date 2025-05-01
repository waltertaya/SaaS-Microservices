package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func reverseProxy(target string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		remote, err := url.Parse(target)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"error": "Bad upstream",
			})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		ctx.Request.URL.Path = ctx.Param("proxyPath") // strip /api/v2/auth or /api/v2/billing

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func main() {
	r := gin.Default()

	// Route /api/v2/auth/* to auth-service
	r.Any("/api/v2/auth/*proxyPath", reverseProxy("http://auth-service:8080"))

	// Route /api/v2/billing/* to billing-service
	r.Any("/api/v2/billing/*proxyPath", reverseProxy("http://billing-service:8081"))

	// API Gateway entrypoint
	r.Run(":8088")
}
