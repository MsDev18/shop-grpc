package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	MySQLDriver = "mysql"
)

type Config struct {
	Host         string `koanf:"host"`
	Port         int    `koanf:"port"`
	Username     string `koanf:"username"`
	Password     string `koanf:"password"`
	Database     string `koanf:"database"`
	MaxOpenConns int    `koanf:"max_open_conns"`
	MaxLifeTime  int    `koanf:"max_life_time"`
	MaxIdleConns int    `koanf:"max_idle_conns"`
}

func (c Config) GetDSN() string {
	DSN := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Database)
	return DSN
}

type Connection struct {
	DB *sql.DB
}

func New(config Config) Connection {
	db, err := sql.Open(MySQLDriver, config.GetDSN())
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * time.Duration(config.MaxLifeTime))
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return Connection{
		DB: db,
	}
}
