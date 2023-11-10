package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/bidon-io/bidon-backend/cmd/bidon-migrate/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./cmd/bidon-migrate/migrations"
	tableName     = "bidon_migrations"

	usagePrefix = `Usage: bidon-migrate [OPTIONS] COMMAND

Options:
`
	usageCommands = `
Commands:
    create NAME [go|sql] Creates new migration file with the current timestamp
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
`
)

func usage() {
	_, _ = fmt.Fprint(os.Stderr, usagePrefix)
	flag.PrintDefaults()
	_, _ = fmt.Fprint(os.Stderr, usageCommands)
}

var (
	verbose = flag.Bool("v", false, "enable verbose mode")

	//go:embed migrations/*.sql
	sqlMigrations embed.FS
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	if *verbose {
		goose.SetVerbose(true)
	}

	switch args[0] {
	case "create":
		if err := goose.Run(args[0], nil, migrationsDir, args[1:]...); err != nil {
			log.Fatal(err)
		}
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

	goose.SetBaseFS(sqlMigrations)
	goose.SetTableName(tableName)

	err = goose.RunWithOptions(
		args[0],
		db,
		"migrations",
		args[1:],
		goose.WithAllowMissing(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
