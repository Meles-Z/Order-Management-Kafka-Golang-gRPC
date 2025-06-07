package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Auth     AuthSecret   `mapstructure:"AUTH"`
	Server   ServerConfig `mapstructure:"SERVER"`
	DBConfig DBConfig     `mapstructure:"DBCONFIG"`
}

type ServerConfig struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

type DBConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Name     string `mapstructure:"NAME"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	SSLMode  string `mapstructure:"SSLMODE"`
}

type AuthSecret struct {
	Websecret string `mapstructure:"WEBSECRET"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Error reading config file:", err)
			return nil, err
		}
		fmt.Println("No .env file found, using environment variables")
	}
	cfg.Server.Host = viper.GetString("HOST")
	cfg.Server.Port = viper.GetString("PORT")

	cfg.DBConfig.Host = viper.GetString("DB_HOST")
	cfg.DBConfig.Port = viper.GetString("DB_PORT")
	cfg.DBConfig.Name = viper.GetString("DB_NAME")
	cfg.DBConfig.User = viper.GetString("DB_USER")
	cfg.DBConfig.Password = viper.GetString("DB_PASSWORD")
	cfg.DBConfig.SSLMode = viper.GetString("SSLMODE")

	cfg.Auth.Websecret = viper.GetString("WEBSECRET")

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
