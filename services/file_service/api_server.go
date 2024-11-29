package main

import (
	"context"
	"cortex_api/common"
	"cortex_api/database"
	"cortex_api/database/db_gen"
	"cortex_api/services/file_service/controllers/v1"
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
			ServerHeader:  "file_service",
			AppName:       "file_service",
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
	//	KeyGenerator:   api_utils.UUID,
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

	fileStoreController := v1.FileStoreController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("get", common.RouteSecureFile+"/:id<guid>", fileStoreController.DownloadFileById)                 // secure download only
	apiServer.Server.Add("get", common.RouteSecureDownload+"/:filename", fileStoreController.DownloadSecureFileByFilename) // secure download only
	apiServer.Server.Add("get", common.RouteImages+"/:filename", fileStoreController.DownloadFileByFilename)               // public access
	apiServer.Server.Add("get", common.RouteDownload+"/:filename", fileStoreController.DownloadFileByFilename)             // public access

	apiServer.Server.Add("get", common.RouteFileStore+"/:id<guid>", fileStoreController.GetFile)
	apiServer.Server.Add("get", common.RouteFileStore+"/filename/:filename", fileStoreController.GetFileByFilename)
	apiServer.Server.Add("get", common.RouteFileStore, fileStoreController.ListFiles)
	apiServer.Server.Add("post", common.RouteFileStore, fileStoreController.UploadFile)
	apiServer.Server.Add("delete", common.RouteFileStore+"/:id<guid>", fileStoreController.DeleteFile)

	return apiServer.Server.Listen(fmt.Sprintf(":%s", os.Getenv("ENDPOINT_PORT")))
}

func (apiServer *Model) StopServer() error {
	return apiServer.Server.Shutdown()
}
