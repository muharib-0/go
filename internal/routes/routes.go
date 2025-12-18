package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muharib-0/ainyx-user-api/internal/handler"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUserByID)
	users.Get("/", userHandler.ListUsers)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
