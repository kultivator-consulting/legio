package context_request

import (
	"cortex_api/common"
	"github.com/gofiber/fiber/v2"
	"os"
)

func New(config Config) fiber.Handler {

	cfg := configDefault(config)

	return func(ctx *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(ctx) {
			return ctx.Next()
		}

		token := common.ExtractToken(ctx)
		if token == "" {
			return ctx.Next()
		}

		publicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
		sub, clientId, err := common.ValidateToken(token, publicKey)
		if err != nil {
			return ctx.Next()
		}

		ctx.Locals(common.ContextToken, token)
		ctx.Locals(common.ContextAccount, sub)
		ctx.Locals(common.ContextSession, clientId)
		return ctx.Next()
	}
}
