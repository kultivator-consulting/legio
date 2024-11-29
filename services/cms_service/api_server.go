package main

import (
	"context"
	"cortex_api/common"
	"cortex_api/database"
	"cortex_api/database/db_gen"
	"cortex_api/services/cms_service/controllers/v1"
	"cortex_api/services/cms_service/interceptors/context_request"
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
			ServerHeader:  "cms_service",
			AppName:       "cms_service",
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

	apiServer.Server.Use(
		context_request.New(context_request.Config{}),
	)

	componentController := v1.ComponentController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("get", common.RouteComponent+"/:id<guid>", componentController.GetComponent)
	apiServer.Server.Add("get", common.RouteComponent, componentController.ListComponents)
	apiServer.Server.Add("post", common.RouteComponent, componentController.CreateComponent)
	apiServer.Server.Add("put", common.RouteComponent+"/:id<guid>", componentController.UpdateComponent)
	apiServer.Server.Add("delete", common.RouteComponent+"/:id<guid>", componentController.DeleteComponent)

	apiServer.Server.Add("post", common.RouteComponent+"/:componentId<guid>/field", componentController.AddComponentField)
	apiServer.Server.Add("delete", common.RouteComponent+"/:componentId<guid>/field/:id<guid>", componentController.RemoveComponentField)

	contentController := v1.ContentController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("get", common.RoutePreviewPage+"/:id<guid>", contentController.GetPageContentById)
	apiServer.Server.Add("get", common.RoutePageContent+"/*", contentController.GetPageContentByFilename)

	apiServer.Server.Add("get", common.RouteContent+"/:id<guid>", contentController.GetContent)
	apiServer.Server.Add("get", common.RouteContent, contentController.ListContents)
	apiServer.Server.Add("post", common.RouteContent, contentController.CreateContent)
	apiServer.Server.Add("post", common.RouteContent+"/:id<guid>", contentController.CloneContent)
	apiServer.Server.Add("put", common.RouteContent+"/:id<guid>", contentController.UpdateContent)
	apiServer.Server.Add("delete", common.RouteContent+"/:id<guid>", contentController.DeleteContent)

	apiServer.Server.Add("post", common.RouteContent+"/:parentId<guid>/child", contentController.AddChildContent)
	apiServer.Server.Add("put", common.RouteContent+"/:parentId<guid>/child/:id<guid>", contentController.UpdateChildContent)
	apiServer.Server.Add("delete", common.RouteContent+"/child/:id<guid>", contentController.RemoveChildContent)

	pagePathController := v1.PagePathController(apiServer.Queries, apiServer.Context)
	apiServer.Server.Add("get", common.RouteBlogs, pagePathController.ListBlogs)

	apiServer.Server.Add("post", common.RoutePagePath+"/:parentId<guid>", pagePathController.CreatePagePath)
	apiServer.Server.Add("post", common.RoutePagePath, pagePathController.CreatePagePath)
	apiServer.Server.Add("put", common.RoutePagePath+"/:id<guid>", pagePathController.UpdatePagePath)
	apiServer.Server.Add("delete", common.RoutePagePath+"/:id<guid>", pagePathController.DeletePagePath)
	apiServer.Server.Add("get", common.RoutePagePath+"/:id<guid>", pagePathController.GetPagePath)
	apiServer.Server.Add("get", common.RoutePagePath, pagePathController.ListPagePaths)
	apiServer.Server.Add("get", common.RoutePagePath+"/links", pagePathController.ListLinks)

	apiServer.Server.Add("post", common.RoutePagePath+"/:parentId<guid>/page", pagePathController.AddPage)
	apiServer.Server.Add("put", common.RoutePagePath+"/page/:id<guid>", pagePathController.UpdatePage)
	apiServer.Server.Add("delete", common.RoutePagePath+"/page/:id<guid>", pagePathController.RemovePage)
	apiServer.Server.Add("get", common.RoutePagePath+"/page/:id<guid>", pagePathController.GetPage)
	apiServer.Server.Add("get", common.RoutePagePath+"/validate/:id<guid>/:slug<string>", pagePathController.ValidatePagePath)
	apiServer.Server.Add("get", common.RoutePagePath+"/draft-page/:id<guid>", pagePathController.CreateDraftPage)

	apiServer.Server.Add("post", common.RoutePagePath+"/blog", pagePathController.AddBlog)
	apiServer.Server.Add("put", common.RoutePagePath+"/blog/:id<guid>", pagePathController.UpdateBlog)
	apiServer.Server.Add("delete", common.RoutePagePath+"/blog/:id<guid>", pagePathController.RemoveBlog)
	apiServer.Server.Add("get", common.RoutePagePath+"/blog/:id<guid>", pagePathController.GetBlog)
	apiServer.Server.Add("get", common.RoutePagePath+"/blog-page/:id<guid>", pagePathController.GetBlogByPage)

	apiServer.Server.Add("post", common.RoutePagePath+"/:parentId<guid>/template", pagePathController.AddPageTemplate)
	apiServer.Server.Add("post", common.RoutePagePath+"/template", pagePathController.AddPageTemplate)
	apiServer.Server.Add("put", common.RoutePagePath+"/template/:id<guid>", pagePathController.UpdatePageTemplate)
	apiServer.Server.Add("delete", common.RoutePagePath+"/template/:id<guid>", pagePathController.RemovePageTemplate)
	apiServer.Server.Add("get", common.RoutePagePath+"/template/:id<guid>", pagePathController.GetPageTemplate)
	apiServer.Server.Add("get", common.RoutePagePath+"/template", pagePathController.ListPageTemplates)

	return apiServer.Server.Listen(fmt.Sprintf(":%s", os.Getenv("ENDPOINT_PORT")))
}

func (apiServer *Model) StopServer() error {
	return apiServer.Server.Shutdown()
}
