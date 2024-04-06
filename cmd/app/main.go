package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"

	"github.com/Kiruhanchik/TestovoyeGo/internal/services"

	"github.com/Kiruhanchik/TestovoyeGo/internal/config"

	"github.com/gorilla/mux"
)

func runMigrations() {
	slog.Info("start migration")

	dbConnStr := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", config.Cfg.User, config.Cfg.Password, config.Cfg.Host, config.Cfg.DbPort, config.Cfg.DbName)

	slog.Debug("conn str " + dbConnStr)

	db, err := sql.Open("postgres", dbConnStr)

	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err)
	}

	err = goose.Up(db, config.Cfg.MigrationDir)
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err)
	}

	slog.Info("complete migration")
}

func main() {
	switch config.Cfg.LogLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	}

	runMigrations()

	router := mux.NewRouter()

	router.HandleFunc("/cars", services.GetCars).Methods("GET")
	router.HandleFunc("/cars/{id}", services.DeleteCar).Methods("DELETE")
	router.HandleFunc("/cars/{id}", services.UpdateCar).Methods("PATCH")
	router.HandleFunc("/cars", services.AddCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+config.Cfg.HTTPAddr, router))
}
