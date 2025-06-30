package main

import (
	"botvideosaver/config"
	tgbot "botvideosaver/internal/bot"
	"botvideosaver/internal/client/pgxpool"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/storage"
	"context"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

//go:generate sqlc generate
func main() {
	logger.InitLogger()

	// // Migrate the database
	// if err := storage.Migrate(); err != nil {
	// 	logger.Log.Sugar().Fatalf("Failed to migrate database: %v", err)
	// }

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().Redis.Host + ":" + config.GetConfig().Redis.Port,
		Password: config.GetConfig().Redis.Password,
		DB:       config.GetConfig().Redis.DB,
	})

	cacheManager := marshaler.New(cache.New[any](redis_store.NewRedis(
		redisClient,
		store.WithExpiration(15*time.Minute),
		store.WithClientSideCaching(5*time.Minute)),
	))

	store := storage.NewStorage(db)

	b, err := tgbot.New(store, cacheManager)
	if err != nil {
		panic(err)
	}

	b.Start(context.Background())
}
