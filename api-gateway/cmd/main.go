package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/waltertaya/saas-microservices/api-gateway/middlewares"
)

func initConfig() {
	viper.SetConfigFile("./config/config.yaml")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

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
	initConfig()
	r := gin.Default()

	// Route /api/v2/auth/* to auth-service
	r.Any("/api/v2/auth/*proxyPath", reverseProxy("http://auth-service:8080"))

	// billing (with JWT)
	billingGroup := r.Group("/api/v2/billing")
	billingGroup.Use(middlewares.JWTMiddleware())
	// Route /api/v2/billing/* to billing-service
	billingGroup.Any("/*proxyPath", reverseProxy("http://billing-service:8081"))

	// API Gateway entrypoint
	r.Run(":8088")
}
