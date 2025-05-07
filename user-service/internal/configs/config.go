package configs

import "github.com/spf13/viper"

type Config struct {
	Auth     AuthConfig `mapstructure:"auth"`
	DBConfig DBConfig   `mapstructure:"dbconfig"`
}

type AuthConfig struct {
	AuthWebSecret string `mapstructure:"AUTH_WEB_SECRET"`
	AuthSecret    string `mapstructure:"AUTH_SECRET"`
}

type DBConfig struct {
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
}

type ServerConfig struct {
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort int    `mapstructure:"SERVER_PORT"`
}

// LoadConfig loads the configuration settings.
// It reads the configuration file, sets up environment variable support, and unmarshal the settings into a Config struct.
// It returns the loaded Config and any error encountered during the process.
func LoadConfig() (*Config, error) {
	var err error
	var config Config

	// AddConfigPath adds the directory where the configuration file is located.
	viper.AddConfigPath(".")

	// SetConfigName sets the name of the configuration file to be read.
	viper.SetConfigName("dev")

	// SetConfigType sets the type of the configuration file.
	viper.SetConfigType("env")

	// AutomaticEnv enables automatic binding of environment variables to configuration values.
	viper.AutomaticEnv()

	// ReadInConfig reads the configuration file with the specified name and type.
	err = viper.ReadInConfig()

	// Check if there was an error reading the configuration file.
	if err != nil {
		return nil, err
	}

	// Unmarshal reads the configuration settings into the Config struct.
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
