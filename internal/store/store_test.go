package store_test

import (
	"log"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/store"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error

	err = godotenv.Load("../../.env.test")
	if err != nil {
		log.Printf("Did not load from .env.test file: %v", err)
	}

	db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = store.AutoMigrate(db)
	if err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}

	os.Exit(m.Run())
}

func ptr[T any](t T) *T {
	return &t
}
