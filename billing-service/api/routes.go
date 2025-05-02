package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/billing-service/controllers"
	"github.com/waltertaya/saas-microservices/billing-service/db"
	"github.com/waltertaya/saas-microservices/billing-service/middlewares"
	"github.com/waltertaya/saas-microservices/billing-service/models"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// protected routes
	protected := r.Group("")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/transactions", func(c *gin.Context) {
		var transactions []models.Transaction
		err := db.DB.Select(&transactions, "SELECT * FROM transactions")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
			return
		}
		c.JSON(http.StatusOK, transactions)
	})
	// simple ping routes for testing
	protected.GET("/welcome", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome you have reached the billing service endpoints",
		})
	})
	// protected.POST("/subscribe", controllers.Subscribe)
	// protected.GET("/subscriptions", controllers.GetSubscriptions)

	// Unprotected routes
	r.POST("/subscribe", controllers.Subscribe)
	r.GET("/subscriptions", controllers.GetSubscriptions)

	r.POST("/pay", controllers.InitiatePayment)
	r.GET("/verify", controllers.VerifyPayment)
	r.POST("/webhook", controllers.HandleWebhook)

	return r
}
