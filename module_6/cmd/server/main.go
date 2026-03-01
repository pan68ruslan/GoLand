package main

import (
	"context"
	"log"
	_ "module_6/cmd/server/docs"
	"module_6/cmd/server/handlers"
	"module_6/cmd/server/middlewares"
	"module_6/internal/clients"
	"module_6/internal/config"
	"module_6/internal/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		//panic(err)
		log.Fatal(err)
	}
	log.Printf("ENV: %s", cfg.Env)
	log.Printf("Mongo URI: %s", cfg.MongoConnectionURI)
	log.Printf("Mongo DB: %s", cfg.MongoDbName)
	log.Printf("Log Level: %s", cfg.LogLevel)
	log.Printf("Token TTL: %s", cfg.TokenTTL)
	log.Printf("Token Secret: %s", cfg.TokenSecret)

	h := handlers.NewHandlers(cfg)
	clients.InitMongo(cfg.MongoConnectionURI)
	defer func() {
		if clients.Mongo != nil {
			if err := clients.Mongo.Disconnect(context.Background()); err != nil {
				utils.Logger.Println("Disconnect Mongo error:", err)
			}
		}
	}()
	app := fiber.New()
	h.SetupRoutes(app)
	// Swagger
	app.Post("/swagger/*", swaggo.HandlerDefault)
	// Auth
	app.Post("/auth/sign-up", handlers.SignUp)
	app.Post("/auth/sign-in", handlers.SignIn)
	// Channel
	app.Post("/channel/history", middlewares.JWTProtected(), handlers.ChannelHistory)
	app.Post("/channel/send", middlewares.JWTProtected(), handlers.ChannelSend)
	// Listen server
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Lunch server error: %v", err)
		}
	}()
	// graceful shutdown (Ctrl+C, Docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Logger.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		utils.Logger.Println("Shutdown error:", err)
	}
}
