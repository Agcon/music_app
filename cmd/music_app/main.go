package main

import (
	"log"
	"music_app/internal/app"
	"music_app/internal/config"
	"music_app/internal/databases"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := databases.NewPostgresClient(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	redisClient := databases.NewRedisClient(cfg.Redis)

	mongoClient, err := databases.NewMongoClient(cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	application := app.NewApp(cfg, db, redisClient, mongoClient)

	if err = application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}

}
