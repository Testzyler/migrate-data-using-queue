package server

import (
	"asynq-quickstart/server/routes"
	"asynq-quickstart/task"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Task payload for any employee related tasks.
type EmployeeTaskPayload struct {
	// ID for the email recipient.
	EmpID int
}

// server api
func Run() {
	// Initialize the server.
	redisAddress := "redis:6379"

	// Initialize the queue client.
	task.Init(redisAddress)

	// Initialize the server.
	app := fiber.New()
	routes.Setup(app)

	// defer task.Close()
	log.Fatal(app.Listen(":3000"))
}
