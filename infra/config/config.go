package config

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host                string `env:"DATABASE_HOST,required"`
	User                string `env:"DATABASE_USER,required"`
	Password            string `env:"DATABASE_PASSWORD,required"`
	Name                string `env:"DATABASE_NAME,required"`
	Port                int    `env:"DATABASE_PORT,required"`
	MaxConnections      int    `env:"DATABASE_MAX_CONNECTIONS,required"`
	MaxIdleConnections  int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" required:"true"`
	DefaultQueryTimeout int    `env:"DATABASE_DEFAULT_QUERY_TIMEOUT_SECS" required:"true"`
}

type Config struct {
	AppServerPort string `env:"PORT,required"`
	Database      DatabaseConfig
}

func extractCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(currentDir + "/.env"); err != nil {
		tokens := strings.Split(currentDir, "/")
		var buf bytes.Buffer
		for i := 0; i < len(tokens)-1; i++ {
			buf.WriteString(tokens[i])
			if i != len(tokens)-2 {
				buf.WriteString("/")
			}
		}
		return buf.String()
	}
	return currentDir
}

func Load() Config {
	currentDir := extractCurrentDir()
	err := godotenv.Load(currentDir + "/.env")
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	log.Default().SetPrefix("\r")

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("unable to parse .env file: %e", err)
	}
	return cfg
}
