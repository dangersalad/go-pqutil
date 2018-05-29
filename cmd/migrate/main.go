package main

import (
	"flag"
	"fmt"
	database "github.com/dangersalad/go-pqutil"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"os"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "/sql", "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	if len(args) > 1 && args[0] == "create" {
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			die(err)
		}
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	if err := goose.SetDialect("postgres"); err != nil {
		die(err)
	}

	db, err := database.Connect(5)
	if err != nil {
		die(err)
	}

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		die(err)
	}
}

func usage() {
	fmt.Print(usagePrefix)
	flags.PrintDefaults()
	fmt.Print(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] COMMAND
Examples:
    migrate status
    migrate create init sql
    migrate create add_some_column sql
    migrate create fetch_user_data go
    migrate up
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
)

func die(err error) {
	fmt.Fprintf(os.Stderr, "%+v", err)
	os.Exit(1)
}
