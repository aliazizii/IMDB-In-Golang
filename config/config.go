package config

import (
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/spf13/viper"
)

type Config struct {
	DB     *Database
	Admin  *model.User
	Secret string `mapstructure:"SECRET"`
}

type Database struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Name     string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

func Read() (*Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	db := &Database{
		Host:     viper.GetString("HOST"),
		Port:     viper.GetString("PORT"),
		Name:     viper.GetString("DB_NAME"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
	}

	admin := &model.User{
		Username: viper.GetString("ADMIN_USERNAME"),
		Password: viper.GetString("ADMIN_PASSWORD"),
	}

	cfg := &Config{
		Secret: viper.GetString("SECRET"),
		DB:     db,
		Admin:  admin,
	}

	return cfg, nil
}
