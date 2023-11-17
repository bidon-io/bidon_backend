// Package dbtest provides helper functions for tests that require database access
package dbtest

import (
	"log"
	"os"

	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bwmarrin/snowflake"
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

	node, err := snowflake.NewNode(0)
	if err != nil {
		log.Fatalf("Error creating snowflake node: %v", err)
	}

	testDB, err = db.Open(os.Getenv("DATABASE_URL"), db.WithSnowflakeNode(node))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = testDB.AutoMigrate()
	if err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}

	return testDB
}
