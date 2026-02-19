package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "192.168.1.20",
		Port:     5432,
		Username: "postgres",
		Password: "password",
		Database: "microservices",
	}
}

func (dbs *DatabaseConfig) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		dbs.Host,
		dbs.Port,
		dbs.Username,
		dbs.Password,
		dbs.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}
