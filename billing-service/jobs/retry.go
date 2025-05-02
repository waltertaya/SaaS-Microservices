// jobs/retry.go
package jobs

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/waltertaya/saas-microservices/billing-service/db"
	"github.com/waltertaya/saas-microservices/billing-service/models"
)

func StartRetryWorker() {
	go func() {
		for {
			var failed []models.FailedWebhook
			db.DB.Select(&failed, "SELECT * FROM failed_webhooks WHERE retry_count < ?", 3)

			for _, f := range failed {
				var payload map[string]interface{}
				json.Unmarshal([]byte(f.Payload), &payload)

				if err := ProcessWebhook(payload); err == nil {
					db.DB.Exec("DELETE FROM failed_webhooks WHERE id = ?", f.ID)
				} else {
					f.RetryCount++
					db.DB.Exec("UPDATE failed_webhooks SET retry_count = ? WHERE id = ?", f.RetryCount, f.ID)
				}
			}

			time.Sleep(30 * time.Second)
		}
	}()
}

func ProcessWebhook(payload map[string]interface{}) error {
	event := payload["event"].(string)
	if event != "charge.completed" {
		return nil // ignore unsupported events
	}

	data := payload["data"].(map[string]interface{})
	status := fmt.Sprintf("%v", data["status"])
	if status != "successful" {
		return fmt.Errorf("transaction not successful")
	}

	// save transaction to DB
	transaction := models.Transaction{
		TransactionID: fmt.Sprintf("%v", data["id"]),
		TxRef:         fmt.Sprintf("%v", data["tx_ref"]),
		Amount:        fmt.Sprintf("%v", data["amount"]),
		Currency:      fmt.Sprintf("%v", data["currency"]),
		Status:        status,
		CustomerEmail: fmt.Sprintf("%v", data["customer"].(map[string]interface{})["email"]),
	}
	query := `INSERT INTO transactions (transaction_id, tx_ref, amount, currency, status, customer_email) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.DB.Exec(query, transaction.TransactionID, transaction.TxRef, transaction.Amount, transaction.Currency, transaction.Status, transaction.CustomerEmail)
	if err != nil {
		return fmt.Errorf("db save failed: %v", err)
	}

	return nil
}
