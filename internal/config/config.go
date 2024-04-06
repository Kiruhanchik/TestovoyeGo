package config

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

var Cfg Config

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8888"`

	// Postgres
	Host         string `env:"POSTGRES_HOST"`
	User         string `env:"POSTGRES_USER"`
	Password     string `env:"POSTGRES_PASSWORD"`
	DbName       string `env:"POSTGRES_DB"`
	DbPort       int    `env:"POSTGRES_PORT"`
	MigrationDir string `env:"MIGRATION_DIR"`
	CarApi       string `env:"CAR_API"`
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("No .env file found")
	}

	if err := envconfig.Process(context.Background(), &Cfg); err != nil {
		log.Fatal(err)
	}
}
