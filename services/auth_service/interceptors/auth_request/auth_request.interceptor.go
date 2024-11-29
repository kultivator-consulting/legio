package auth_request

import (
	"cortex_api/common"
	"cortex_api/database/db_gen"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
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
			log.Printf("AuthRequestInterceptor invalid raw access token\n")
			return cfg.Unauthorized(ctx)
		}

		publicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
		sub, clientId, err := common.ValidateToken(token, publicKey)
		if err != nil {
			log.Printf("AuthRequestInterceptor error validating token, %v\n", err)
			return cfg.Unauthorized(ctx)
		}

		sessions, err := cfg.Queries.GetSessionsByAccountIdClientId(cfg.DbCtx, db_gen.GetSessionsByAccountIdClientIdParams{
			AccountID: sub.(pgtype.UUID),
			ClientID:  clientId.(pgtype.UUID),
		})
		if err != nil {
			log.Printf("AuthRequestInterceptor error while validating session, error %v\n", err)
			return cfg.Unauthorized(ctx)
		}
		if len(sessions) != 1 {
			// no matching session found
			return cfg.Unauthorized(ctx)
		}

		ctx.Locals(common.ContextToken, token)
		ctx.Locals(common.ContextAccount, sub)
		ctx.Locals(common.ContextSession, clientId)
		return ctx.Next()
	}
}
