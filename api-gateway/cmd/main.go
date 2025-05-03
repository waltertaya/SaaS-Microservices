package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/waltertaya/saas-microservices/api-gateway/middlewares"
	"github.com/waltertaya/saas-microservices/api-gateway/utils"
	"go.uber.org/zap"
)

func initConfig() {
	viper.SetConfigFile("./config/config.yaml")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

type contextKey string

const traceIDKey contextKey = "trace_id"

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

		// Inject trace ID into the request context for upstream & error handling
		traceID, exists := ctx.Get("trace_id")
		if exists {
			ctx.Request = ctx.Request.WithContext(
				context.WithValue(ctx.Request.Context(), traceIDKey, traceID),
			)
		}

		// Custom error handler to log upstream errors
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			traceVal := req.Context().Value("trace_id")
			traceID := "unknown"
			if traceVal != nil {
				traceID = traceVal.(string)
			}

			utils.Logger.Error("proxy error",
				zap.String("trace_id", traceID),
				zap.String("target", target),
				zap.Error(err),
			)

			rw.WriteHeader(http.StatusBadGateway)
			rw.Write([]byte("Upstream service error"))
		}

		// Modify the request path (retain my original logic) : bug
		ctx.Request.URL.Path = ctx.Param("proxyPath")

		// Forward to upstream
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func main() {
	initConfig()
	utils.InitLogger()
	r := gin.Default()

	r.Use(gin.Recovery())        // error recovery
	r.Use(utils.RequestLogger()) // Add logging

	// simple ping route
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Route /api/v2/auth/* to auth-service
	// r.Any("/api/v2/auth/*proxyPath", reverseProxy("http://auth-service:8080"))
	r.Any("/api/v2/auth/*proxyPath", reverseProxy("http://localhost:8080"))

	// billing (with JWT)
	billingGroup := r.Group("/api/v2/billing")
	billingGroup.Use(middlewares.JWTMiddleware())
	// Route /api/v2/billing/* to billing-service
	// billingGroup.Any("/*proxyPath", reverseProxy("http://billing-service:8081"))
	billingGroup.Any("/*proxyPath", reverseProxy("http://localhost:8081"))

	// API Gateway entrypoint
	r.Run(":8088")
}
