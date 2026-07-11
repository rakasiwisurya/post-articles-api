package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"

	"post-articles-api/internal/article"
	"post-articles-api/internal/config"
	"post-articles-api/internal/database"
)

func main() {
	cfg := config.Load()

	if err := database.EnsureDatabase(cfg); err != nil {
		log.Fatalf("ensure database: %v", err)
	}

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db, cfg.DBName); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	app := fiber.New(fiber.Config{AppName: "post-articles-api"})
	app.Use(recoverer.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		ExposeHeaders: []string{"X-Total-Count"},
	}))

	handler := article.NewHandler(article.NewService(article.NewRepository(db)))
	article.RegisterRoutes(app, handler)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
