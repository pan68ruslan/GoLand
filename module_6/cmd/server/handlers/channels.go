package handlers

import (
	"context"
	"time"

	"module_6/internal/clients"
	"module_6/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// Message — модель повідомлення
type Message struct {
	Username  string    `json:"username" bson:"username"`
	Content   string    `json:"content" bson:"content"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// ChannelHistory — повертає останні повідомлення з каналу
func ChannelHistory(c fiber.Ctx) error {
	collection := clients.Mongo.Database("chatdb").Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}, nil)
	if err != nil {
		utils.Logger.Println("history error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch history"})
	}
	defer cursor.Close(ctx)

	var messages []Message
	if err := cursor.All(ctx, &messages); err != nil {
		utils.Logger.Println("parsing history:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not parse history"})
	}

	return c.JSON(messages)
}

// ChannelSend — додає нове повідомлення у канал
func ChannelSend(c fiber.Ctx) error {
	collection := clients.Mongo.Database("chatdb").Collection("messages")

	var msg Message
	if err := c.Bind().Body(&msg); err != nil {
		utils.Logger.Println("parsing history:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	msg.Timestamp = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, msg)
	if err != nil {
		utils.Logger.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not save message"})
	}

	utils.Logger.Println("new messages:", msg.Username)
	return c.JSON(fiber.Map{"status": "message sent"})
}
