package utils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./db/otp_store.db"

// GenerateOTP generates a 6-digit OTP and stores it in the database with a 10-minute expiration.
func GenerateOTP(username string) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Create table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS otps (
		username TEXT,
		otp TEXT,
		expiration_time DATETIME
	);`
	if _, err := db.Exec(createTable); err != nil {
		return "", err
	}

	// Generate OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	expiration := time.Now().Add(10 * time.Minute)

	// Insert OTP
	insertStmt := `INSERT INTO otps (username, otp, expiration_time) VALUES (?, ?, ?);`
	_, err = db.Exec(insertStmt, username, otp, expiration.Format(time.RFC3339))
	if err != nil {
		return "", err
	}

	return otp, nil
}

// ValidateOTP checks if the given OTP for the username is valid and not expired.
func ValidateOTP(username, otp string) (bool, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := `SELECT expiration_time FROM otps WHERE username = ? AND otp = ? LIMIT 1;`
	var expirationStr string
	err = db.QueryRow(query, username, otp).Scan(&expirationStr)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	expiration, err := time.Parse(time.RFC3339, expirationStr)
	if err != nil {
		return false, err
	}

	if time.Now().After(expiration) {
		return false, nil
	}

	return true, nil
}

// ClearExpiredOTPs deletes OTP records that have expired.
func ClearExpiredOTPs() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM otps WHERE expiration_time < ?`, time.Now().Format(time.RFC3339))
	if err != nil {
		fmt.Println("Error clearing expired OTPs:", err)
	}
}
