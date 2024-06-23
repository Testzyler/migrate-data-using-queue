package app

import (
	"asynq-quickstart/task"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetEmployeeToMigrate(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Employee data"})
}

func RemoveEmployee(c *fiber.Ctx) error {
	type request struct {
		EmpID string `json:"emp_id"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	empID := req.EmpID
	migrator := task.EmployeeMigrator{}
	tasks, err := migrator.NewRemoveEmployeeTasks([]string{empID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create migrate employee tasks"})
	}
	client := task.GetClient()
	// defer client.Close()

	for _, tk := range tasks {
		info, err := client.Enqueue(tk)
		if err != nil {
			fmt.Printf("Error enqueuing task: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not enqueue image resize task"})
		}

		fmt.Printf("Task enqueued: %v\n", info)
	}

	return c.JSON(fiber.Map{"message": "Employee data removed"})
}
