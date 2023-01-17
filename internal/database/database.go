package database

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/config"
)

var (
	conn *connectionOptions
)

// connectionOptions describes necessary information
// to access a Postgres instance
type connectionOptions struct {
	host     string
	port     int64
	database string
	username string
	password string
}

// String converts connection options to a format
// that the database library understands
func (c *connectionOptions) string() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.username, c.password, c.database,
	)
}

// Init stores the database connection options
func Init() {
	conn = &connectionOptions{
		config.C.PostgresHost,
		config.C.PostgresPort,
		config.C.PostgresDatabase,
		config.C.PostgresUsername,
		config.C.PostgresPassword,
	}
}

// Connect tries to connect to the database
// using options that describe the necessary
// information
func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", conn.string())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Ping tries to connect to the database
// and pings this to test if the connection is okay
func Ping() error {
	db, err := Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

// Exec executes a single query in the database
func Exec(query string, params ...interface{}) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	res, err := db.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("unable to execute query '%s' with params '%v': %v", query, params, err)
	}

	if _, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	return nil
}
