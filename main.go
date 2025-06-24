package main

import (
	"botvideosaver/config"
	tgbot "botvideosaver/internal/bot"
	"botvideosaver/internal/client/pgxpool"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/storage"
	"context"
)

func main() {
	logger.InitLogger()

	db, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		Url:             config.GetConfig().Postgres.Url,
		Host:            config.GetConfig().Postgres.Host,
		Port:            config.GetConfig().Postgres.Port,
		Username:        config.GetConfig().Postgres.Username,
		Password:        config.GetConfig().Postgres.Password,
		Database:        config.GetConfig().Postgres.Database,
		MaxConnections:  config.GetConfig().Postgres.MaxConnections,
		MaxConnIdleTime: config.GetConfig().Postgres.MaxConnIdleTime,
	})
	if err != nil {
		logger.Log.Sugar().Errorf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		logger.Log.Sugar().Errorf("Failed to ping database: %v", err)
	}

	logger.Log.Sugar().Infof("Connected to database %s at %s:%d",
		config.GetConfig().Postgres.Database,
		config.GetConfig().Postgres.Host,
		config.GetConfig().Postgres.Port,
	)

	store := storage.NewStorage(db)

	b, err := tgbot.New(store)
	if err != nil {
		panic(err)
	}

	b.Start(context.Background())
}
