package models

import "time"

// NotificationChannels represents user's preferred notification methods
type NotificationChannels struct {
	Email             string `json:"email" bson:"email"`
	SMS               string `json:"sms" bson:"sms"`
	PushNotifications bool   `json:"push_notifications" bson:"push_notifications"`
}

// Subscription represents user's topic subscriptions
type Subscription struct {
	UserID               string               `json:"user_id" bson:"user_id"`
	Topics               []string             `json:"topics" bson:"topics"`
	NotificationChannels NotificationChannels `json:"notification_channels" bson:"notification_channels"`
	CreatedAt            time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" bson:"updated_at"`
}

// EventDetails represents the details of a notification event
type EventDetails struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Message represents the notification message content
type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// NotificationEvent represents a complete notification event
type NotificationEvent struct {
	Topic   string       `json:"topic"`
	EventID string       `json:"event_id"`
	Time    time.Time    `json:"timestamp"`
	Details EventDetails `json:"details"`
	Message Message      `json:"message"`
}

// UnsubscribeRequest represents the unsubscribe payload
type UnsubscribeRequest struct {
	UserID string   `json:"user_id"`
	Topics []string `json:"topics"`
}
