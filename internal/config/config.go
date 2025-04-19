package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	cfg := &Config{}

	appCfg, err := loadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}
	cfg.App = appCfg

	postgresCfg, err := loadPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load postgres config: %w", err)
	}
	cfg.Postgres = postgresCfg

	mongoCfg, err := loadMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load mongo config: %w", err)
	}
	cfg.Mongo = mongoCfg

	redisCfg, err := loadRedisConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load redis config: %w", err)
	}
	cfg.Redis = redisCfg

	return cfg, nil
}

func loadAppConfig() (*AppConfig, error) {
	httpPort := os.Getenv("HTTP_PORT")
	jwtKey := os.Getenv("JWT_SECRET")
	jwtTTl := os.Getenv("JWT_TTL")
	postgresConfigFileName := os.Getenv("POSTGRES_CONFIG_PATH")
	mongoConfigFileName := os.Getenv("MONGO_CONFIG_PATH")
	redisConfigFileName := os.Getenv("REDIS_CONFIG_PATH")
	if httpPort == "" || jwtKey == "" || jwtTTl == "" || postgresConfigFileName == "" || mongoConfigFileName == "" || redisConfigFileName == "" {
		return nil, fmt.Errorf("requirable environment variables is not set")
	}
	intJwtTTl, _ := strconv.Atoi(jwtTTl)
	return &AppConfig{
		HTTPPort:               httpPort,
		JWTKey:                 jwtKey,
		JWTTTL:                 intJwtTTl,
		PostgresConfigFileName: postgresConfigFileName,
		MongoConfigFileName:    mongoConfigFileName,
		RedisConfigFileName:    redisConfigFileName,
	}, nil
}

func loadPostgresConfig() (*PostgresConfig, error) {
	file, err := getConfigFile("POSTGRES_CONFIG_PATH")
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres config file: %w", err)
	}
	defer file.Close()

	postgresCfg := &PostgresConfig{}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&postgresCfg); err != nil {
		return nil, fmt.Errorf("failed to decode postgres config file: %w", err)
	}
	return postgresCfg, nil
}

func loadMongoConfig() (*MongoConfig, error) {
	file, err := getConfigFile("MONGO_CONFIG_PATH")
	if err != nil {
		return nil, fmt.Errorf("failed to open mongo config file: %w", err)
	}
	defer file.Close()

	mongoCfg := &MongoConfig{}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&mongoCfg); err != nil {
		return nil, fmt.Errorf("failed to decode mongo config file: %w", err)
	}
	return mongoCfg, nil
}

func loadRedisConfig() (*RedisConfig, error) {
	file, err := getConfigFile("REDIS_CONFIG_PATH")
	if err != nil {
		return nil, fmt.Errorf("failed to open redis config file: %w", err)
	}
	defer file.Close()

	redisCfg := &RedisConfig{}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&redisCfg); err != nil {
		return nil, fmt.Errorf("failed to decode redis config file: %w", err)
	}
	return redisCfg, nil
}

func getConfigFile(envName string) (*os.File, error) {
	configPath := os.Getenv(envName)
	if configPath == "" {
		return nil, fmt.Errorf("requirable environment variables is not set")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open mongo config file: %w", err)
	}
	return file, nil
}
