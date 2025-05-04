package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
	conf Config
)

type Config struct {
	AppPort string `envconfig:"app_port"`

	DBName    string `envconfig:"db_name"`
	DBUser    string `envconfig:"db_user"`
	DBPass    string `envconfig:"db_pass"`
	DBHost    string `envconfig:"db_host"`
	DBPort    string `envconfig:"db_port"`
	DBMaxIdle int    `envconfig:"db_max_idle"`
	DBMaxOpen int    `envconfig:"db_max_open"`
}

func GetConfig() Config {
	once.Do(func() {
		_ = godotenv.Load()
		if err := envconfig.Process("", &conf); err != nil {
			log.Fatal(err)
		}
	})
	return conf
}
