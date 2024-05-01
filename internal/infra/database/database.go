package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var instance *sql.DB
var once sync.Once

type DBPool struct {
	dsn                string
	maxConnections     int
	maxIdleConnections int
}
type DBConfig struct {
	Host               string
	User               string
	Password           string
	DBName             string
	Port               int
	MaxConnections     int
	MaxIdleConnections int
}

func NewDBPool(cfg DBConfig) DBPool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.DBName)

	return DBPool{dsn: dsn, maxConnections: cfg.MaxConnections,
		maxIdleConnections: cfg.MaxIdleConnections}
}

func (p *DBPool) GetDB() *sql.DB {
	once.Do(func() {
		db, err := sql.Open("mysql", p.dsn)
		if err != nil {
			panic(err.Error())
		}
		db.SetMaxIdleConns(p.maxIdleConnections)
		db.SetMaxOpenConns(p.maxConnections)
		instance = db
	})
	return instance
}

func (p *DBPool) GetDSN() string {
	return p.dsn
}
