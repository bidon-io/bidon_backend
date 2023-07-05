// Package dbtest provides helper functions for tests that require database access
package dbtest

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/joho/godotenv"
)

func Prepare() *db.DB {
	var (
		testDB *db.DB
		err    error
	)

	err = godotenv.Load("../../../.env.test")
	if err != nil {
		log.Printf("Did not load from .env.test file: %v", err)
	}

	testDB, err = db.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = testDB.AutoMigrate()
	if err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}

	return testDB
}

func CreateUser(t *testing.T, tx *db.DB, index int) *db.User {
	t.Helper()

	user := &db.User{
		Email: fmt.Sprintf("test%d@email.com", index),
	}

	if err := tx.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	return user
}

func CreateUsersList(t *testing.T, tx *db.DB, usersCount int) []*db.User {
	t.Helper()

	users := make([]*db.User, usersCount)
	for i := range users {
		users[i] = CreateUser(t, tx, i)
	}

	return users
}
