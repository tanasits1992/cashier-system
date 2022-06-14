package utils

import (
	"net/http"

	"github.com/gofiber/fiber"
)

func BadRequestResponse(c *fiber.Ctx, message string) {
	c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"error": true,
		"msg":   message,
	})
}

func InternalServerErrorResponse(c *fiber.Ctx, message string) {
	c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": true,
		"msg":   message,
	})
}

func SuccessResponse(c *fiber.Ctx, data interface{}) {
	c.JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  data,
	})
}
