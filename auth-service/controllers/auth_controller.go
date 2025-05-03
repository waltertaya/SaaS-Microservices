package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/auth-service/db"
	"github.com/waltertaya/saas-microservices/auth-service/models"
	"github.com/waltertaya/saas-microservices/auth-service/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user.Password = string(hashedPassword)
	user.Role = "user"

	query := `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`

	err = db.DB.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User

	query := `SELECT * FROM users WHERE username=$1`
	err := db.DB.Get(&user, query, input.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// forgot password and send email
func ForgotPassword(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
	}
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	query := `SELECT * FROM users WHERE username=$1`
	err = db.DB.Get(&user, query, input.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Username not found",
		})
		return
	}

	var otp string
	otp, err = utils.GenerateOTP(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate OTP",
		})
		return
	}
	err = utils.SendEmail(user.Email, otp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send email",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to your email",
	})
}

// verify OTP
func VerifyOTP(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		OTP      string `json:"otp"`
	}
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	valid, err := utils.ValidateOTP(input.Username, input.OTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate OTP",
		})
		return
	}
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid OTP",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
	})
}

// reset password
func ResetPassword(ctx *gin.Context) {
	var input struct {
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	query := `UPDATE users SET password=$1 WHERE username=$2`
	_, err = db.DB.Exec(query, hashedPassword, input.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update password",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
