package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	AllowPreviousPublications bool
	PostgresPort int64
	PostgresUsername string
	PostgresPassword string
	PostgresHost string
	PostgresDatabase string
	BindAddress string
	PublisherName string
	PublisherEmail string
	ApplicationDomain string
	SubscriptionAESPassword string
	SendGridAPIKey string
	MinScrapeInterval time.Duration
}

var C *config

var Defaults = map[string]interface{}{
	"POSTGRES_USERNAME": "newsletter",
	"POSTGRES_PASSWORD": "newsletter",
	"POSTGRES_HOST":     "localhost",
	"POSTGRES_PORT":     5432,
	"POSTGRES_DATABASE": "newsletter",
	"BIND_ADDRESS":      "127.0.0.1:8080",
	"ALLOW_PREVIOUS_PUBLICATIONS": "true",
	"SENDGRID_API_KEY": "CHANGEME",
	"PUBLISHER_EMAIL": "text@example.com",
	"PUBLISHER_NAME": "Example User",
	"SUBSCRIPTION_AES_PASSWORD": "CHANGEME",
	"MIN_SCRAPE_INTERVAL": "168h",  // 7 days
	"APPLICATION_DOMAIN": "newsletter.statictask.io",
}

func Initialize() {
	C = &config{
		AllowPreviousPublications: getEnvOrDefaultBool("ALLOW_PREVIOUS_PUBLICATIONS"),
		PostgresPort: getEnvOrDefaultInt64("POSTGRES_PORT"),
		PostgresUsername: getEnvOrDefaultString("POSTGRES_USERNAME"),
		PostgresPassword: getEnvOrDefaultString("POSTGRES_PASSWORD"),
		PostgresHost: getEnvOrDefaultString("POSTGRES_HOST"),
		PostgresDatabase: getEnvOrDefaultString("POSTGRES_DATABASE"),
		BindAddress: getEnvOrDefaultString("BIND_ADDRESS"),
		PublisherName: getEnvOrDefaultString("PUBLISHER_NAME"),
		PublisherEmail: getEnvOrDefaultString("PUBLISHER_EMAIL"),
		ApplicationDomain: getEnvOrDefaultString("APPLICATION_DOMAIN"),
		SendGridAPIKey: getEnvOrDefaultString("SENDGRID_API_KEY"),
		MinScrapeInterval: getEnvOrDefaultDuration("MIN_SCRAPE_INTERVAL"),
	}
}

// getEnvOrDefault returns the environment variable
// value returned by viper or the hardcoded default
func getEnvOrDefault(key string) interface{} {
	v := viper.Get(strings.ToLower(key))
	if v == nil {
		return Defaults[key] 
	}

	return v
}

// getEnvOrDefaultInt64 returns the environment variable
// value returned by viper or the hardcoded default
func getEnvOrDefaultInt64(key string) int64 {
	var err error
	value, ok := getEnvOrDefault(key).(int)
	if !ok {
		str := getEnvOrDefault(key).(string)	
		value, err = strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
	}

	return int64(value)
}

// getEnvOrDefaultBool returns the environment variable
// value returned by viper or the hardcoded default
func getEnvOrDefaultBool(key string) bool {
	value := getEnvOrDefault(key).(string)
	return strings.ToLower(value) == "true"
}

// getEnvOrDefaultString returns the environment variable
// value returned by viper or the hardcoded default
func getEnvOrDefaultString(key string) string {
	return getEnvOrDefault(key).(string)
}

// getEnvOrDefaultDuration returns the environment variable
// value returned by viper or the hardcoded default
func getEnvOrDefaultDuration(key string) time.Duration {
	duration, err := time.ParseDuration(getEnvOrDefault(key).(string))
	if err != nil {
		panic(err)
	}

	return duration
}
