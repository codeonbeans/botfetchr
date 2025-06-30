package storage

import (
	"botvideosaver/config"
	"botvideosaver/internal/client/pgxpool"
	"botvideosaver/internal/logger"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

const migrationsDir = "migrations"

func Migrate() error {
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		logger.Log.Fatal("Failed to find migrations directory", zap.Error(err))
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Log.Fatal("Failed to set dialect", zap.Error(err))
		return err
	}

	c := config.GetConfig().Postgres
	connectionString := pgxpool.GetConnStr(pgxpool.PgxpoolOptions{
		Url:             c.Url,
		Host:            c.Host,
		Port:            c.Port,
		Username:        c.Username,
		Password:        c.Password,
		Database:        c.Database,
		MaxConnections:  c.MaxConnections,
		MaxConnIdleTime: c.MaxConnIdleTime,
	})

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Log.Error("Failed to open database connection", zap.Error(err))
		return err
	}

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.Error("Failed to close database connection", zap.Error(err))
		}
	}(db)

	if err = db.Ping(); err != nil {
		return err
	}

	logger.Log.Info("Database connection established")
	if err = goose.Up(db, migrationsDir); err != nil {
		logger.Log.Error("Failed to run migrations", zap.Error(err))
		return err
	}

	return nil
}
