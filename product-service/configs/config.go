package configs

import (
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
	viper.SetConfigName("env") // file name without extension
	viper.SetConfigType("env") // type of file
	viper.AddConfigPath("./")  // look for the file in this path
	viper.AutomaticEnv()       // automatically override with env variables if present

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
