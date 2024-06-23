package database

import (
	"asynq-quickstart/env"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

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

func (d *Database) Migrate() {
	d.DB.AutoMigrate()
}

// Get DB instance
func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

func (d *Database) Close() {
	sqlDB, err := d.DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}
