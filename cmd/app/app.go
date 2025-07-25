package main

import (
	"github.com/Cherillkin/Notepad/config"
	"github.com/Cherillkin/Notepad/database"
	"github.com/Cherillkin/Notepad/handlers"
	"github.com/Cherillkin/Notepad/repositories"
	"github.com/Cherillkin/Notepad/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	envConfig := config.NewEnvConfig()
	db := database.Init(envConfig, database.DBMigrator)

	app := fiber.New(fiber.Config{
		AppName:      "NotePad",
		ServerHeader: "Fiber",
	})

	// listRepository := repositories.NewListRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	authService := services.NewAuthService(authRepository)

	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	// privateRoutes := app.Use(middlewares.AuthProtected(db))

	app.Listen(":8000")
}
