package db

import (
	"fmt"
	"log"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	connectionstring = "host=localhost user=jamboree password=61c9b888a55fdbd3f18bc6d6 dbname=jamboree sslmode=disable"
)

var db *gorm.DB

func init() {
	fmt.Println("database initialized")
	connectToDB()
	runAutoMigrations()
}

func runAutoMigrations() {
	db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Address{},
	)
}

func connectToDB() {
	var err error
	db, err = gorm.Open(postgres.Open(connectionstring), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
}
