package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/waltertaya/saas-microservices/billing-service/api"
	"github.com/waltertaya/saas-microservices/billing-service/db"
	"github.com/waltertaya/saas-microservices/billing-service/jobs"
)

func initConfig() {
	viper.SetConfigFile("./config/config.yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initConfig()

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	fmt.Println("Billing DB connected ✅")

	jobs.StartRetryWorker()

	r := api.SetupRouter()

	port := viper.GetString("server.port")
	if port == "" {
		port = "8081" // fallback default
	}

	r.Run(":" + port)
}
