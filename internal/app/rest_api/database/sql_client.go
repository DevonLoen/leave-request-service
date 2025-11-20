package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	DBDriver          string
	DBSource          string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxIdleTime   time.Duration
	ConnectionTimeout time.Duration
}

type SQLClient struct {
	DB *sql.DB
}

func NewSQLClient(cfg Config) (*SQLClient, error) {
	db, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &SQLClient{DB: db}, nil
}

func (client *SQLClient) Close() error {
	if client.DB != nil {
		return client.DB.Close()
	}

	return nil
}
