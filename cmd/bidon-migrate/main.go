package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/bidon-io/bidon-backend/cmd/bidon-migrate/migrations"
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/db/gen"
	_ "github.com/jackc/pgx/v5/stdlib"
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
    generate-models      Generates models from the DB schema
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
	noGen   = flag.Bool("no-gen", false, "disable model generation")

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

	config.LoadEnvFile()
	if config.GetEnv() == config.UnknownEnv {
		slog.Error("ENVIRONMENT is required")
		os.Exit(1)
	}

	if config.GetEnv() == config.ProdEnv && os.Getenv("IKNOWWHATIAMDOING") != "yes" {
		if args[0] != "status" && args[0] != "version" && args[0] != "up" {
			slog.Error("only 'status', 'version' and 'up' commands are allowed in production environment. Use 'IKNOWWHATIAMDOING=yes' to override.")
			os.Exit(1)
		}
	}

	if *verbose {
		goose.SetVerbose(true)
	}

	// Handle commands that do not need DB connection
	switch args[0] {
	case "create":
		if err := goose.Run(args[0], nil, migrationsDir, args[1:]...); err != nil {
			slog.Error("goose.Run()", "error", err)
			os.Exit(1)
		}
		return
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("missing DATABASE_URL environment variable")
		os.Exit(1)
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		slog.Error("failed to open DB", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close DB", "error", err)
			os.Exit(1)
		}
	}()

	// Handle custom commands that need DB connection
	switch args[0] {
	case "generate-models":
		gen.GenerateModels(db)
		return
	}

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
		slog.Error("goose.RunWithOptions()", "error", err)
		os.Exit(1)
	}
	// Generate models after migration tasks in dev environment
	if config.GetEnv() == config.DevEnv && !*noGen {
		gen.GenerateModels(db)
	}
}
