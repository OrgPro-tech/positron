package config

import "github.com/spf13/viper"

type Config struct {
	DB_Url        string `mapstructure:"db_url"`
	ServerPort    string `mapstructure:"Server_port"`
	Database_Name string `mapstructure:"Database_name"`
	Username      string `mapstructure:"Username"`
	Password      string `mapstructure:"Password"`
	Host          string `mapstructure:"Host"`
	DB_Port       string `mapstructure:"DB_Port"`
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
