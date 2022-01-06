package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/database"
	"fiber-mongo-api/models"
	"fiber-mongo-api/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	// Creates a new Fiber instance.
	app := fiber.New(fiber.Config{
		AppName:      configs.Getenv("APP_NAME"),
		ServerHeader: configs.Getenv("APP_HEADER"),
	})

	// Connect to database
	database.ConnectDB()
	defer database.Cancel()
	defer database.Client.Disconnect(database.Ctx)

	// Create model schemas
	models.CreateUserSchema()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	//routes
	routes.SetupRoutes(app)

	app.Listen(fmt.Sprintf(":%s", configs.Getenv("APP_PORT")))
}
