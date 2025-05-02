package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/waltertaya/saas-microservices/billing-service/db"
	"github.com/waltertaya/saas-microservices/billing-service/models"
)

// should be created immediately user is registered
func Subscribe(ctx *gin.Context) {
	var input struct {
		Plan string `json:"plan"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := ctx.GetInt("user_id")

	sub := models.Subscription{
		UserID:    userID,
		Plan:      input.Plan,
		Status:    "active",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 1, 0), // one month from now
	}

	query := `INSERT INTO subscriptions (user_id, plan, status, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.DB.QueryRow(
		query, sub.UserID, sub.Plan, sub.Status, sub.StartDate, sub.EndDate,
	).Scan(&sub.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "Subscribed successfully",
		"subscription": sub,
	})
}

func GetSubscriptions(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	var subs []models.Subscription

	query := `SELECT * FROM subscriptions WHERE user_id=$1`

	err := db.DB.Select(&subs, query, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscriptions": subs,
	})
}
