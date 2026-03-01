package handlers

import (
	"module_6/cmd/server/handlers/health"
	"module_6/internal/config"
	//"module_6/cmd/server/handlers/"

	"github.com/gofiber/fiber/v3"
)

type Handlers struct {
	//cfg *config.Config
	//Env string
	Health *health.Handler `json:"Health,omitempty"`
}

func NewHandlers(cfg *config.Config) *Handlers {
	return &Handlers{
		//cfg: cfg,
		//Env: cfg.Env,
	}
}

func (h *Handlers) SetupRoutes(app *fiber.App) {
	app.Get("/health", h.Health.Health)
	/*app.Get("/swagger/*", swaggerui.New(swaggerui.Config{
		BasePath: "/swagger",
		FilePath: "/docs/swagger.json",
	}))

	api := app.Group("/api/v1")
	api.Use(h.middlewares.Log.Handle)

	counterGroup := api.Group("/counter")
	counterGroup.Get("/increment", h.Counter.Increment)

	authGroup := api.Group("/auth")
	authGroup.Post("/signup", h.Auth.SignUp)
	authGroup.Post("/signin", h.Auth.SignIn)

	messageGroup := api.Group("/message")
	messageGroup.Use(h.middlewares.Auth.Handle)
	messageGroup.Post("/send", h.Message.SendMessage)
	messageGroup.Post("/list", h.Message.GetMessageList)
	messageGroup.Get("/listen/:channel_name", websocket.New(h.Message.Listen))*/
}
