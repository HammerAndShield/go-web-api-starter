package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-web-api-starter/internal/common"
	"os"
	"time"
)

type SqlDbConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type DB struct {
	*sql.DB
}

func NewDatabase(sslEnabled bool) *SqlDbConfig {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	var dsn string

	if sslEnabled {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	}

	return &SqlDbConfig{
		Dsn:          dsn,
		MaxOpenConns: common.IntEnv(os.Getenv, "DB_MAX_OPEN_CONNS", 100),
		MaxIdleConns: common.IntEnv(os.Getenv, "DB_MAX_IDLE_CONNS", 50),
		MaxIdleTime:  30 * time.Minute,
	}
}

func (sqlDbConfig SqlDbConfig) OpenDB(driver string) (*DB, error) {
	db, err := sql.Open(driver, sqlDbConfig.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(sqlDbConfig.MaxOpenConns)
	db.SetMaxIdleConns(sqlDbConfig.MaxIdleConns)
	db.SetConnMaxIdleTime(sqlDbConfig.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}
