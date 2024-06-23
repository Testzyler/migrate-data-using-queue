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
	redisAddress := "localhost:6379"
	app := fiber.New()
	routes.Setup(app)

	// Initialize the queue client.
	task.Init(redisAddress)
	// defer task.Close()
	log.Fatal(app.Listen(":3000"))
}
