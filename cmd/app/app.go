package main

import (
	"github.com/Cherillkin/Notepad/config"
	"github.com/Cherillkin/Notepad/database"
	"github.com/Cherillkin/Notepad/handlers"
	"github.com/Cherillkin/Notepad/middlewares"
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

	listRepository := repositories.NewListRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	authService := services.NewAuthService(authRepository)
	listService := services.NewListService(listRepository)

	server := app.Group("/api")

	authGroup := server.Group("/auth")
	handlers.NewAuthHandler(authGroup, authService, db)

	privateRoutes := server.Group("/list", middlewares.AuthProtected(db))
	handlers.NewListHandler(privateRoutes, listService)

	app.Listen(":8000")
}
