package main

import (
	"context"
	"cortex_api/common"
	"cortex_api/database"
	"cortex_api/database/db_gen"
	"cortex_api/services/auth_service/controllers/v1"
	"cortex_api/services/auth_service/interceptors/auth_request"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"log"
	"os"
	"time"
)

type Model struct {
	Server  *fiber.App
	Pool    *pgxpool.Pool
	Queries *db_gen.Queries
	Context context.Context
}

type Interface interface {
	StartServer() error
	StopServer() error
}

var Module = fx.Options(fx.Provide(ApiServer))

func ApiServer() *Model {
	return &Model{
		Server: fiber.New(fiber.Config{
			StrictRouting: true,
			ServerHeader:  "auth_service",
			AppName:       "auth_service",
		}),
	}
}

func (apiServer *Model) StartServer() error {
	pool, queries, ctx, err := database.ApiDatabase().Open()
	if err != nil {
		log.Fatalf("StartServer DB connection error: %v\n", err)
	}
	apiServer.Pool = pool
	apiServer.Queries = queries
	apiServer.Context = ctx

	defer func(database *database.Model) {
		err := database.Close()
		if err != nil {
			log.Fatalf("StartServer DB close error: %v\n", err)
		}
	}(database.ApiDatabase())

	apiServer.Server.Use(
		helmet.New(),
	)

	// TODO - add CSRF protection
	//csrfConfig := csrf.Config{
	//	KeyLookup:      "header:X-Xsrf-Token",
	//	CookieName:     "doug_dug_",
	//	CookieSameSite: "Strict",
	//	Expiration:     3 * time.Hour,
	//	KeyGenerator:   utils.UUID,
	//}
	//apiServer.Server.Use(
	//	csrf.New(csrfConfig),
	//)

	limiterConfig := limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        10000,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			// TODO - Implement error and static status pages
			log.Printf("Limiter threshold reached\n")
			return c.SendFile("./too-fast-page.html")
		},
	}
	apiServer.Server.Use(
		limiter.New(limiterConfig),
	)

	apiServer.Server.Use(
		logger.New(logger.Config{
			TimeFormat: time.RFC3339,
			TimeZone:   "UTC",
		}),
	)

	authController := v1.AuthController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("post", common.RouteAuth, authController.AuthLogin)
	apiServer.Server.Add("post", common.RouteAuth+"/forgot-password", authController.AuthForgotPassword)
	apiServer.Server.Add("post", common.RouteAuth+"/reset-password", authController.AuthResetPassword)
	apiServer.Server.Add("get", common.RouteAuth, authController.AuthRefresh)

	// NOTE: Order here is important, it means that the auth_request middleware will be called after
	// auth_controller but before all other secure controllers.
	authRequestConfig := auth_request.Config{
		Queries: apiServer.Queries,
		DbCtx:   apiServer.Context,
	}
	apiServer.Server.Use(
		auth_request.New(authRequestConfig),
	)

	apiServer.Server.Add("get", common.RouteToken, authController.AuthValidate)

	apiServer.Server.Add("head", common.RouteAuth, authController.AuthValidate)
	apiServer.Server.Add("delete", common.RouteAuth, authController.AuthLogout)

	accountController := v1.AccountController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("get", common.RouteAccount+"/:id", accountController.GetAccount)
	apiServer.Server.Add("get", common.RouteAccount, accountController.ListAccounts)
	apiServer.Server.Add("post", common.RouteAccount, accountController.CreateAccount)
	apiServer.Server.Add("put", common.RouteAccount+"/:id", accountController.UpdateAccount)
	apiServer.Server.Add("delete", common.RouteAccount+"/:id", accountController.DeleteAccount)

	profileController := v1.ProfileController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("get", common.RouteProfile, profileController.GetProfile)
	apiServer.Server.Add("put", common.RouteProfile, profileController.UpdateProfile)
	apiServer.Server.Add("delete", common.RouteProfile, profileController.DeleteProfile)

	return apiServer.Server.Listen(fmt.Sprintf(":%s", os.Getenv("ENDPOINT_PORT")))
}

func (apiServer *Model) StopServer() error {
	return apiServer.Server.Shutdown()
}
