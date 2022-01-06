package routes

import (
	"fiber-mongo-api/controllers"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/", hello)

	users := apiv1.Group("/users")
	users.Get("/", controllers.GetUsers)
	users.Get("/:id", controllers.GetUser)
	users.Post("/", controllers.CreateUser)
	users.Put("/:id", controllers.UpdateUser)
	users.Delete("/:id", controllers.DeleteUser)

}

func hello(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(&fiber.Map{"message": "API is running"})
}
