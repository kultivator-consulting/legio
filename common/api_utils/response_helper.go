package api_utils

import "github.com/gofiber/fiber/v2"

func SetRESTHeaders(ctx *fiber.Ctx) {
	ctx.Set("Content-Type", "charset=utf-8, application/json")
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
	ctx.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Set("Access-Control-Expose-Headers", "Set-Cookie")
}
