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
	itemRepository := repositories.NewItemRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	listService := services.NewListService(listRepository)
	itemService := services.NewItemService(itemRepository)
	authService := services.NewAuthService(authRepository)

	server := app.Group("/api")

	privateRoutes := server.Group("/list", middlewares.AuthProtected(db))
	handlers.NewListHandler(privateRoutes, listService)

	itemRoutes := server.Group("/:listId/item", middlewares.AuthProtected(db), middlewares.SetListIdToLocals)
	handlers.NewItemHandler(itemRoutes, itemService)

	authGroup := server.Group("/auth")
	handlers.NewAuthHandler(authGroup, authService, db)

	app.Listen(":8000")
}
