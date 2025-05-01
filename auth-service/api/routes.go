package api

import (
	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/auth-service/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/v2/auth/register", controllers.Register)
	r.POST("/api/v2/auth/login", controllers.Login)

	return r
}
