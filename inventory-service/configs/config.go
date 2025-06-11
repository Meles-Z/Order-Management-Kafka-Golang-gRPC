package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Auth Auth
	DB   DatabaseConfig
	ENV  Env
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type Auth struct {
	Secret string
}

type Env struct {
	Env string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Error reading config file:", err)
			return nil, err
		}
		fmt.Println("No .env file found, using system env vars")
	}

	cfg := &Config{
		Auth: Auth{
			Secret: viper.GetString("SECRET"),
		},
		DB: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		ENV: Env{
			Env: viper.GetString("ENV"),
		},
	}

	return cfg, nil
}