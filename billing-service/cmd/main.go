package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/waltertaya/saas-microservices/billing-service/api"
	"github.com/waltertaya/saas-microservices/billing-service/db"
)

func main() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	fmt.Println("Billing DB connected ✅")

	r := api.SetupRouter()
	port := viper.GetString("server.port")
	r.Run(":" + port)
}
