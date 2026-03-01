package health

import (
	"encoding/json"
	"module_6/internal/config"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	env string
}

type Response struct {
	Env string `json:"env"`
	Ok  bool   `json:"ok"`
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		env: cfg.Env,
	}
}

func (h *Handler) Health(ctx fiber.Ctx) error {
	res, err := json.Marshal(Response{
		Env: h.env,
		Ok:  true,
	})
	if err != nil {
		return err
	}
	return ctx.Send(res)
}
