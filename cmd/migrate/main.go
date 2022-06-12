//go:build migrate
// +build migrate

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	cmd := flag.String("cmd", "up", "[up or down]")
	flag.Parse()

	pgUsername := getEnv("NEWSLETTER_POSTGRES_USERNAME", "newsletter")
	pgPassword := getEnv("NEWSLETTER_POSTGRES_PASSWORD", "newsletter")
	pgDatabase := getEnv("NEWSLETTER_POSTGRES_DATABASE", "newsletter")
	pgHost := getEnv("NEWSLETTER_POSTGRES_HOST", "localhost")
	pgPort := getEnv("NEWSLETTER_POSTGRES_PORT", "5432")

	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			pgUsername,
			pgPassword,
			pgHost,
			pgPort,
			pgDatabase,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	switch string(*cmd) {
	case "up":
		log.Println("Starting migration: UP")
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		log.Println("Starting migration: DOWN")
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("command not found")
	}

	log.Println("Finished")
}
