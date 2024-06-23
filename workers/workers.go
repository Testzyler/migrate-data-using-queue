package workers

import (
	"asynq-quickstart/database"
	"asynq-quickstart/task"
	"context"
	"log"

	"github.com/hibiken/asynq"
)

const redisAddr = "localhost:6379"

func Run() {
	ctx := context.Background()
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"export":  5,
				"purge":   4,
				"default": 1,
			},
			BaseContext: func() context.Context { return ctx },
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	// Register tasks with their corresponding handlers
	handler := task.EmployeeMigrator{
		Ctx: &ctx,
		DB:  db,
	}

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeMigrateEmployee, handler.HandleMigrateEmployeeTask)
	mux.HandleFunc(task.TypeRemoveEmployee, handler.HandleRemoveEmployeeTask)
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
