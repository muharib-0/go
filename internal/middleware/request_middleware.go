package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/muharib-0/ainyx-user-api/internal/logger"
	"go.uber.org/zap"
)

func RequestID() fiber.Handler {
	return requestid.New()
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		logger.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("request_id", c.Get(fiber.HeaderXRequestID)),
			zap.String("ip", c.IP()),
		)

		return err
	}
}
