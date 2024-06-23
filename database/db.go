package database

import (
	"asynq-quickstart/env"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	host := env.GetEnv("Database.Host")
	port := env.GetEnv("Database.Port")
	user := env.GetEnv("Database.Username")
	password := env.GetEnv("Database.Password")
	dbname := env.GetEnv("Database.Name")
	schema := env.GetEnv("Database.Schema")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host, port, user, password, dbname, schema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        1000,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to database [%s] on %s:%s\n", dbname, host, port)
	return db, nil
}
