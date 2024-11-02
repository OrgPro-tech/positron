package config

import "github.com/spf13/viper"

type Config struct {
	DB_Url        string `mapstructure:"DATABASE_URL"`
	ServerPort    string `mapstructure:"Server_port"`
	Database_Name string `mapstructure:"DATABASE_NAME"`
	Username      string `mapstructure:"DATABASE_USERNAME"`
	Password      string `mapstructure:"DATABASE_PASSWORD"`
	Host          string `mapstructure:"DATABASE_HOST"`
	DB_Port       string `mapstructure:"DATABASE_PORT"`
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil
	}

	return &config
}
