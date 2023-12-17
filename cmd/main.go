package main

import (
	"log/slog"

	"github.com/lbsti/eulabs-challenge/adapter/api"
	migrate "github.com/lbsti/eulabs-challenge/db"
	"github.com/lbsti/eulabs-challenge/infra/config"
	"github.com/lbsti/eulabs-challenge/infra/database"
	"github.com/lbsti/eulabs-challenge/infra/repository"
)

func main() {
	cfg := config.Load()
	dbPool := database.NewDBPool(database.DBConfig{
		Host:               cfg.Database.Host,
		Port:               cfg.Database.Port,
		User:               cfg.Database.User,
		Password:           cfg.Database.Password,
		DBName:             cfg.Database.Name,
		MaxConnections:     cfg.Database.MaxConnections,
		MaxIdleConnections: cfg.Database.MaxIdleConnections,
	})
	db := dbPool.GetDB()
	defer db.Close()

	e := migrate.RunMigrate("up", dbPool.GetDSN(), "db/migrations", []string{}...)
	if e != nil {
		slog.Error("failure when execute migration %v", e)
	}
	productRepo := repository.NewProductRepositorySQL(db)
	webServer := api.NewWebServer(cfg.AppServerPort, productRepo)
	webServer.Run()
}
