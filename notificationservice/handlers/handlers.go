package Handlers

import (
	"net/http"
	"notificationservice/models"
	Services "notificationservice/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *Services.NotificationService
}

func NewNotificationHandler(ns *Services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: ns,
	}
}

func (h *NotificationHandler) Subscribe(c *gin.Context) {
	var subscription models.Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.Subscribe(&subscription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully subscribed", "subscription": subscription})
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var event models.NotificationEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.SendNotification(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

func (h *NotificationHandler) Unsubscribe(c *gin.Context) {
	var req models.UnsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.Unsubscribe(req.UserID, req.Topics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed"})
}

func (h *NotificationHandler) GetSubscriptions(c *gin.Context) {
	userID := c.Param("user_id")

	subscriptions, err := h.notificationService.GetSubscriptions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
