package routes

import (
	hr "asynq-quickstart/server/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// Register routes
	app.Post("/migrate/employee", hr.GetEmployeeToMigrate)
	app.Post("/remove/employee", hr.RemoveEmployee)
}
