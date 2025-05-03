package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/auth-service/db"
	"github.com/waltertaya/saas-microservices/auth-service/models"
)

// simple protected test route
func Profile(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	role := ctx.GetString("role")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your profile",
		"user_id": userID,
		"role":    role,
	})
}

// update profile
func UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Update details required",
		})
		return
	}

	var user models.User

	err = db.DB.Get(&user, "SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	if input.Username == "" {
		input.Username = user.Username
	}
	if input.Email == "" {
		input.Email = user.Email
	}

	// update user
	_, err = db.DB.Exec("UPDATE users SET username = ?, email = ? WHERE id = ?", input.Username, input.Email, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": input.Username,
			"email":    input.Email,
			"role":     user.Role,
		},
	})
}

// delete profile
func DeleteProfile(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
