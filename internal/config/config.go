package config

import (
	"fmt"
	"log"
	"os"

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
	}

	configName := fmt.Sprintf("config.%s", env)

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Unable to decode config:", err)
	}

	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	if frontend := os.Getenv("FRONTEND_URL"); frontend != "" {
		cfg.Server.FrontendURL = frontend
	}

	return &cfg
}
