package api

import (
	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/auth-service/controllers"
	"github.com/waltertaya/saas-microservices/auth-service/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/forgot-password", controllers.ForgotPassword)
	r.POST("/verify-otp", controllers.VerifyOTP)
	r.POST("/reset-password", controllers.ResetPassword)

	// Test simple protected user profile
	protected := r.Group("")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/profile", controllers.Profile)

	return r
}
