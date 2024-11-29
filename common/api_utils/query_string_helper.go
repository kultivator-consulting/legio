package api_utils

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func ParseQueryString(ctx *fiber.Ctx) (url.Values, error) {
	return url.ParseQuery(ctx.Get("x-query-string"))
}
