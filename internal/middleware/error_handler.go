package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muharib-0/ainyx-user-api/internal/logger"
	"go.uber.org/zap"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		logger.Error("Request error", zap.Error(err))

		return c.Status(code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
}
