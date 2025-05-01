package api

import (
	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/billing-service/controllers"
	"github.com/waltertaya/saas-microservices/billing-service/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	protected := r.Group("/api/v2/billing")
	protected.Use(middlewares.AuthMiddleware())
	protected.POST("/subscribe", controllers.Subscribe)
	protected.GET("subscriptions", controllers.GetSubscriptions)

	return r
}
