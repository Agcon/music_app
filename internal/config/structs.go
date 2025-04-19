package config

import "time"

type Config struct {
	App      *AppConfig
	Postgres *PostgresConfig
	Mongo    *MongoConfig
	Redis    *RedisConfig
}

type PostgresConfig struct {
	Host     string `yaml:"postgres_host"`
	Port     int    `yaml:"postgres_port"`
	User     string `yaml:"postgres_user"`
	Password string `yaml:"postgres_password"`
	DBName   string `yaml:"postgres_dbname"`
	SSLMode  string `yaml:"postgres_sslmode,omitempty"`
}

type MongoConfig struct {
	URI      string `yaml:"mongo_uri"`
	Database string `yaml:"mongo_database"`
}

type RedisConfig struct {
	Addr     string        `yaml:"redis_addr"`
	Password string        `yaml:"redis_password"`
	DB       int           `yaml:"redis_db,omitempty"`
	TTL      time.Duration `yaml:"redis_ttl,omitempty"`
}

type AppConfig struct {
	HTTPPort               string `env:"HTTP_PORT" envDefault:":8080"`
	JWTKey                 string `env:"JWT_SECRET,required"`
	JWTTTL                 int    `env:"JWT_TTL" envDefault:"86400"`
	PostgresConfigFileName string `env:"POSTGRES_CONFIG_PATH" envDefault:"postgres.yaml"`
	MongoConfigFileName    string `env:"MONGO_CONFIG_PATH" envDefault:"mongo.yaml"`
	RedisConfigFileName    string `env:"REDIS_CONFIG_PATH" envDefault:"redis.yaml"`
}
