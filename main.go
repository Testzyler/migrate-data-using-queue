package main

import (
	"asynq-quickstart/server"
	"asynq-quickstart/workers"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [server|worker]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "server":
		server.Run()
	case "worker":
		workers.Run()
	default:
		fmt.Println("Invalid command. Use 'server' or 'worker'.")
		os.Exit(1)
	}
}
