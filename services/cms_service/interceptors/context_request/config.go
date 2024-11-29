package context_request

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Next           func(c *fiber.Ctx) bool
	InvalidContext fiber.Handler
}

var ConfigDefault = Config{
	Next:           nil,
	InvalidContext: nil,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.InvalidContext == nil {
		cfg.InvalidContext = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	return cfg
}
