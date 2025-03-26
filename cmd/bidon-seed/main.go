package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/bidon-io/bidon-backend/config"
)

//go:embed seeds/*.sql
var seedMigrations embed.FS

func main() {
	config.LoadEnvFile()

	if config.GetEnv() == config.TestEnv {
		log.Println("Skipping seeds in test environment.")
		return
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("missing DATABASE_URL environment variable")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("failed to open DB: ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("failed to close DB: ", err)
		}
	}()

	seedsFS, err := fs.Sub(seedMigrations, "seeds")
	if err != nil {
		log.Fatal(err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		seedsFS,
		goose.WithDisableVersioning(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = provider.Up(context.Background()); err != nil {
		log.Fatal(err)
	}
}
