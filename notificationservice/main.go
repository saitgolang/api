package main

import (
	"context"
	"log"
	Handlers "notificationservice/handlers"
	Services "notificationservice/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("notifications")

	// Initialize Kafka service
	kafkaService, err := Services.NewKafkaService([]string{"localhost:9092"})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize services and handlers
	notificationService := Services.NewNotificationService(db, kafkaService)
	notificationHandler := Handlers.NewNotificationHandler(notificationService)

	// Setup Gin router
	router := gin.Default()

	// Routes
	router.POST("/subscribe", notificationHandler.Subscribe)
	router.POST("/notifications/send", notificationHandler.SendNotification)
	router.POST("/unsubscribe", notificationHandler.Unsubscribe)
	router.GET("/subscriptions/:user_id", notificationHandler.GetSubscriptions)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
