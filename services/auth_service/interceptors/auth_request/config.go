package auth_request

import (
	"context"
	"cortex_api/database/db_gen"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Next         func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Queries      *db_gen.Queries
	DbCtx        context.Context
}

var ConfigDefault = Config{
	Next:         nil,
	Unauthorized: nil,
	Queries:      nil,
	DbCtx:        nil,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	if cfg.Queries == nil {
		cfg.Queries = ConfigDefault.Queries
	}
	if cfg.DbCtx == nil {
		cfg.DbCtx = ConfigDefault.DbCtx
	}
	return cfg
}
