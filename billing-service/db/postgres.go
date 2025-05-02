package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sqlx.DB

func Connect() error {
	// dbConfig := viper.Sub("database")

	dbConfig := viper.Sub("database")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.GetString("host"),
		dbConfig.GetString("port"),
		dbConfig.GetString("user"),
		dbConfig.GetString("password"),
		dbConfig.GetString("dbname"),
		dbConfig.GetString("sslmode"),
	)

	// dsn := fmt.Sprintf(
	// 	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	viper.GetString("DB_HOST"),
	// 	viper.GetString("DB_PORT"),
	// 	viper.GetString("DB_USER"),
	// 	viper.GetString("DB_PASSWORD"),
	// 	viper.GetString("DB_NAME"),
	// 	// dbConfig.GetString("sslmode"),
	// )

	var err error
	DB, err = sqlx.Connect("postgres", dsn)

	// var err error
	// for i := 1; i <= 10; i++ {
	// 	DB, err = sqlx.Connect("postgres", dsn)
	// 	if err == nil {
	// 		return nil
	// 	}
	// 	fmt.Printf("DB connect attempt %d failed: %v; retrying in 2s\n", i, err)
	// 	time.Sleep(2 * time.Second)
	// }

	// return fmt.Errorf("could not connect after retries: %w", err)
	return err
}
