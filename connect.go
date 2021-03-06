// Package pqutil has utilities for working with postgres
package pqutil // import "github.com/dangersalad/go-pqutil"

import (
	"database/sql"
	"fmt"
	"time"

	env "github.com/dangersalad/go-environment"
	// we are using postgres for this
	_ "github.com/lib/pq"
)

const (
	// EnvKeyHost is the postgres host (default is "localhost")
	EnvKeyHost = "DB_HOST"
	// EnvKeyPort is the postgres port (default is "5432")
	EnvKeyPort = "DB_PORT"
	// EnvKeyUser is the postgres user
	EnvKeyUser = "DB_USER"
	// EnvKeyPassword is the postgres password
	EnvKeyPassword = "DB_PASSWORD"
	// EnvKeyDatabase is the postgres database to connect to
	EnvKeyDatabase = "DB_DATABASE"
	// EnvKeySSLMode is the postgres ssl mode to use (default is "disable")
	EnvKeySSLMode = "DB_SSL_MODE"
)

func reattemptConnect(attempts int, err error) (*sql.DB, error) {
	attempts--
	if attempts == 0 {
		return nil, err
	}
	debugf("error connecting to database, %d attempts remaining: %s", attempts, err)
	time.Sleep(2 * time.Second)
	return Connect(attempts)
}

var envVars env.Options

// Connect will connect to the database, trying until it connects, or
// the supplied number of attempts have been made to connect. Since
// this runs in k8s, the database proxy container may not be fully
// ready when this attempts to connect for the first time.
func Connect(attempts int) (*sql.DB, error) {

	params, err := env.ReadOptions(env.Options{
		EnvKeyHost:     "localhost",
		EnvKeyPort:     "5432",
		EnvKeyUser:     "",
		EnvKeyPassword: "",
		EnvKeyDatabase: "",
		EnvKeySSLMode:  "disable",
	})

	if err != nil {
		return nil, err
	}

	envVars = params

	return ConnectTo(attempts,
		params[EnvKeyHost],
		params[EnvKeyPort],
		params[EnvKeyUser],
		params[EnvKeyPassword],
		params[EnvKeyDatabase],
		params[EnvKeySSLMode])

}

// ConnectTo connects to the specified database given the supplied credentials
func ConnectTo(attempts int, host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return reattemptConnect(attempts, fmt.Errorf("connecting to database: %w", err))
	}

	err = db.Ping()
	if err != nil {
		return reattemptConnect(attempts, fmt.Errorf("pinging database: %w", err))
	}

	if sslmode == "disable" {
		logf("connected to %s with SSL disabled", host)
	}

	return db, nil
}

// GetUsername will return the username that connected to the database
func GetUsername() string {
	return envVars[EnvKeyUser]
}

// GetDatabaseName will return the database name that was connected to
func GetDatabaseName() string {
	return envVars[EnvKeyDatabase]
}
