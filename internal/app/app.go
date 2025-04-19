package app

import (
	"github.com/redis/go-redis/v9"
	"music_app/internal/config"
	"music_app/internal/databases"
	"music_app/internal/music/handler"
	"music_app/internal/music/repository"
	"music_app/internal/music/service"
	"music_app/internal/user"
	"music_app/pkg/auth"
	"time"
)

type App struct {
	cfg            *config.Config
	pg             *databases.SQLDatabase
	redisClient    *redis.Client
	userRepo       user.Repository
	jwtManager     auth.JWTManager
	sessionManager auth.SessionManager
	mongoClient    *databases.MongoClient
	trackRepo      repository.TrackRepository
	trackService   service.TrackService
	trackHandler   *handler.TrackHandler
}

func NewApp(cfg *config.Config, pg *databases.SQLDatabase, redisClient *redis.Client, mongoClient *databases.MongoClient) *App {
	return &App{
		cfg:            cfg,
		pg:             pg,
		redisClient:    redisClient,
		userRepo:       user.NewRepository(pg.DB),
		jwtManager:     auth.NewJWTManager(cfg.App.JWTKey, intToDuration(cfg.App.JWTTTL)),
		sessionManager: auth.NewRedisSession(redisClient),
		mongoClient:    mongoClient,
	}
}

func intToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
