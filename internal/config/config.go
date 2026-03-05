package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	JWT struct {
		Secret string
	}
}

func LoadConfig() *Config {

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

	return &cfg
}
