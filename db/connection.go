package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=localhost user=gymstats password=123 dbname=gymstatsdb port=7000 "

var DB *gorm.DB

func DBConnection() {
	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Fatal("Error connecting to database")
	} else {
		log.Println("Connected to database")
	}
}
