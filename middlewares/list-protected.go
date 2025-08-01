package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetListIdToLocals(ctx *fiber.Ctx) error {
	listIDUuint, err := strconv.ParseUint(ctx.Params("listId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid list ID",
		})
	}
	ctx.Locals("listId", uint(listIDUuint))
	return ctx.Next()
}
