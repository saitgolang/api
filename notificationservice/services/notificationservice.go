package Services

import (
	"context"
	"notificationservice/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationService struct {
	db           *mongo.Database
	kafkaService *KafkaService
}

func NewNotificationService(db *mongo.Database, kafkaService *KafkaService) *NotificationService {
	return &NotificationService{
		db:           db,
		kafkaService: kafkaService,
	}
}

func (ns *NotificationService) Subscribe(subscription *models.Subscription) error {
	collection := ns.db.Collection("subscriptions")
	ctx := context.Background()

	// Set timestamps
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	// Check if subscription exists
	var existing models.Subscription
	err := collection.FindOne(ctx, bson.M{"user_id": subscription.UserID}).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// Create new subscription
		_, err = collection.InsertOne(ctx, subscription)
		return err
	}

	// Update existing subscription
	update := bson.M{
		"$set": bson.M{
			"topics":                subscription.Topics,
			"notification_channels": subscription.NotificationChannels,
			"updated_at":            time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"user_id": subscription.UserID}, update)
	return err
}

func (ns *NotificationService) Unsubscribe(userID string, topics []string) error {
	collection := ns.db.Collection("subscriptions")
	ctx := context.Background()

	update := bson.M{
		"$pull": bson.M{
			"topics": bson.M{
				"$in": topics,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	return err
}

func (ns *NotificationService) GetSubscriptions(userID string) (*models.Subscription, error) {
	collection := ns.db.Collection("subscriptions")
	ctx := context.Background()

	var subscription models.Subscription
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (ns *NotificationService) SendNotification(event *models.NotificationEvent) error {
	return ns.kafkaService.PublishMessage(event.Topic, event)
}
