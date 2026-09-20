package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	FrontendURL string `mapstructure:"frontend_url"`
}

type DatabaseConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	DBName     string `mapstructure:"dbname"`
	SSLMode    string `mapstructure:"sslmode"`
	AutoCreate bool   `mapstructure:"auto_create"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "dev"
		if os.Getenv("VERCEL") != "" {
			env = "production"
		}
	}

	configName := fmt.Sprintf("config.%s", env)

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	cfg := Config{
		Environment: env,
		Server: ServerConfig{
			Port: "8080",
		},
		Database: DatabaseConfig{
			Port:       5432,
			DBName:     "postgres",
			SSLMode:    "require",
			AutoCreate: false,
		},
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file unavailable, using environment variables: %v", err)
	} else if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("unable to decode config file, using environment variables: %v", err)
	}

	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	if frontend := os.Getenv("FRONTEND_URL"); frontend != "" {
		cfg.Server.FrontendURL = frontend
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}

	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.DBName = dbName
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		port, err := strconv.Atoi(dbPort)
		if err != nil {
			log.Printf("invalid DB_PORT %q, using %d", dbPort, cfg.Database.Port)
		} else {
			cfg.Database.Port = port
		}
	}

	return &cfg
}
