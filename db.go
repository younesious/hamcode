package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	user     = "postgres"
	password = "yourpass"
	host     = "172.20.0.2"
	port     = 5432
	dbname   = "postgres"
)

var doOnce sync.Once
var singleton *gorm.DB

func GetConnection() *gorm.DB {
	doOnce.Do(func() {
		connURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

		db, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // Use standard output for logs
				logger.Config{
					LogLevel: logger.Info, // Set log level to Info to print SQL queries
				},
			),
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		err = db.AutoMigrate(&Attendance{})
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}

		singleton = db
	})
	return singleton
}
