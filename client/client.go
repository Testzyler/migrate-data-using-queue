package client

import (
	"asynq-quickstart/task"
	"log"

	"github.com/hibiken/asynq"
)

// client.go
func Run() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
	empID := "MQ991023-001"
	migrator := task.EmployeeMigrator{}
	tasks, err := migrator.NewRemoveEmployeeTasks([]string{empID})
	if err != nil {
		log.Fatal(err)
	}

	for _, tk := range tasks {
		// Process the task immediately.
		info, err := client.Enqueue(tk)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(" [*] Successfully enqueued task: %+v", info)
	}
}
