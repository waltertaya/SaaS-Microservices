package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", int(claims["user_id"].(float64)))
		ctx.Set("role", claims["role"].(string))
		ctx.Next()
	}
}
