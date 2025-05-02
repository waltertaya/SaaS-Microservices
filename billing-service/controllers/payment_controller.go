package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/waltertaya/saas-microservices/billing-service/db"
	"github.com/waltertaya/saas-microservices/billing-service/jobs"
	"github.com/waltertaya/saas-microservices/billing-service/models"
	"github.com/waltertaya/saas-microservices/billing-service/utils"
)

func InitiatePayment(c *gin.Context) {
	type PaymentRequest struct {
		Amount        string `json:"amount"`
		CustomerEmail string `json:"email"`
	}

	var req PaymentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	req.Amount = "200"

	client := resty.New()

	// Use the environment variables
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing secret key"})
		return
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+secretKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"tx_ref":          "tx-" + utils.GenerateRef(), // implement a proper unique ref
			"amount":          req.Amount,
			"currency":        "NGN",
			"redirect_url":    "https://your-redirect-url.com", // FRONTEND REDIRECT
			"payment_options": "card",
			"customer": map[string]string{
				"email": req.CustomerEmail,
			},
			"customizations": map[string]string{
				"title":       "Test Payment",
				"description": "Pay with card on Flutterwave (Test)",
			},
		}).
		Post("https://api.flutterwave.com/v3/payments")

	if err != nil {
		log.Println("Flutterwave error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment initiation failed"})
		return
	}

	c.Data(resp.StatusCode(), "application/json", resp.Body())
}

func VerifyPayment(c *gin.Context) {
	transactionID := c.Query("transaction_id")
	userID := strconv.Itoa(c.GetInt("user_id"))
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing transaction id"})
		return
	}
	// fmt.Println(transactionID)

	client := resty.New()
	secretKey := os.Getenv("SECRET_KEY")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+secretKey).
		Get("https://api.flutterwave.com/v3/transactions/" + transactionID + "/verify")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "verification failed"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid JSON response"})
		return
	}

	if result["status"] == "success" {
		data := result["data"].(map[string]interface{})
		customer := data["customer"].(map[string]interface{})

		transaction := models.Transaction{
			UserID:        userID,
			TransactionID: transactionID,
			TxRef:         fmt.Sprintf("%v", data["tx_ref"]),
			Amount:        fmt.Sprintf("%v", data["amount"]),
			Currency:      fmt.Sprintf("%v", data["currency"]),
			Status:        fmt.Sprintf("%v", data["status"]),
			CustomerEmail: fmt.Sprintf("%v", customer["email"]),
		}

		// set transaction date to current time
		transaction.TransactionDate = time.Now().Format(time.RFC3339)

		_, err := db.DB.NamedExec(`INSERT INTO transactions (transaction_id, tx_ref, amount, currency, status, customer_email, transaction_date) 
			VALUES (:transaction_id, :tx_ref, :amount, :currency, :status, :customer_email, :transaction_date)`, &transaction)
		if err != nil {
			log.Println("Failed to save transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "Payment verified and saved",
			"transaction": transaction,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "transaction not successful"})
}

func HandleWebhook(c *gin.Context) {
	webhookSecret := os.Getenv("WEBHOOK_SECRET")

	// Validate signature from Flutterwave
	signature := c.GetHeader("verif-hash")
	if signature == "" || signature != webhookSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
		return
	}

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	event, ok := payload["event"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event"})
		return
	}

	// Process the webhook event
	if err := jobs.ProcessWebhook(payload); err != nil {
		// If processing fails, store it for retry
		data, _ := json.Marshal(payload)
		failed := models.FailedWebhook{
			Event:      event,
			Payload:    string(data),
			RetryCount: 0,
		}
		_, err := db.DB.NamedExec(`INSERT INTO failed_webhooks (event, payload, retry_count) VALUES (:event, :payload, :retry_count)`, &failed)
		if err != nil {
			log.Println("Failed to save webhook:", err)
		}

		log.Println("Webhook processing failed. Saved for retry:", err)
	}

	log.Println("✅ Webhook Event:", event)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
