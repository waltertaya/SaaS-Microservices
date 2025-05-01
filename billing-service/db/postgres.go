package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sqlx.DB

func Connect() error {
	dbConfig := viper.Sub("database")

	dsn := fmt.Sprintf(
		"host=%s port=%s password=%s dbname=%s sslmode=%s",
		dbConfig.GetString("host"),
		dbConfig.GetString("port"),
		dbConfig.GetString("password"),
		dbConfig.GetString("dbname"),
		dbConfig.GetString("sslmode"),
	)

	var err error

	DB, err = sqlx.Connect("postgres", dsn)
	return err
}
