package main

import (
	"log"

	"github.com/lbsti/eulabs-challenge/adapter/api"
	migrate "github.com/lbsti/eulabs-challenge/db"
	"github.com/lbsti/eulabs-challenge/internal/infra/config"
	"github.com/lbsti/eulabs-challenge/internal/infra/database"
	"github.com/lbsti/eulabs-challenge/internal/infra/repository"
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

	err := migrate.RunMigrate("up", dbPool.GetDSN(), "db/migrations", []string{}...)
	if err != nil {
		log.Default().Printf("failure when execute migration %v", err)
	}
	productRepo := repository.NewProductRepositorySQL(db)
	webServer := api.NewWebServer(cfg.AppServerPort, productRepo)
	webServer.Run()
}
