package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muharib-0/ainyx-user-api/config"
	sqlc "github.com/muharib-0/ainyx-user-api/db/sqlc"
	"github.com/muharib-0/ainyx-user-api/internal/handler"
	applogger "github.com/muharib-0/ainyx-user-api/internal/logger"
	"github.com/muharib-0/ainyx-user-api/internal/repository"
	"github.com/muharib-0/ainyx-user-api/internal/routes"
	"github.com/muharib-0/ainyx-user-api/internal/service"
)

func main() {
	// Initialize logger
	applogger.Init()
	defer applogger.Sync()

	// Load configuration - using LoadConfig (not Load)
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize SQLC queries
	queries := sqlc.New(pool)

	// Initialize repository
	userRepo := repository.NewUserRepository(queries)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Ainyx User API v1.0.0",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(app, userHandler)

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Starting server on %s", serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
