package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"dev"`
	Port        string `env:"PORT" envDefault:"8080"`
	SwaggerHost string `env:"SWAGGER_HOST" envDefault:"localhost:8080"`
	MongoURI    string `env:"MONGO_URI" envDefault:""`
	MongoDB     string `env:"MONGO_DB" envDefault:"secret-santa"`
}

var Cfg *Config

func LoadConfig() {
	Cfg = new(Config)
	if err := env.Parse(Cfg); err != nil {
		log.Fatalf(`error on parse env variables due:[%v]`, err)
	}
}
