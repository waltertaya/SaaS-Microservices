package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/waltertaya/saas-microservices/auth-service/api"
	"github.com/waltertaya/saas-microservices/auth-service/db"
)

func main() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	fmt.Println("Connected to DB ✅")

	r := api.SetupRouter()

	port := viper.GetString("server.port")

	r.Run(":" + port)
}
