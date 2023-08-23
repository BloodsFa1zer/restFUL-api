package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	URI_BD string `env:"URI_MongoDB"`
	DBName string `env:"DatabaseName"`
	DBPath string `env:"DatabasePath"`
}

func LoadENV(filename string) *Config {
	err := godotenv.Load(filename)
	if err != nil {
		log.Panic().Err(err).Msg(" does not load .env")
	}
	log.Info().Msg("successfully load .env")
	cfg := Config{}
	return &cfg

}

func (cfg *Config) ParseENV() {

	err := env.Parse(cfg)
	if err != nil {
		log.Panic().Err(err).Msg(" unable to parse environment variables")
	}
	log.Info().Msg("successfully parsed .env")
}

func (cfg *Config) MongoENV() string {
	cfg.ParseENV()
	return cfg.URI_BD
}

func (cfg *Config) DBNameENV() string {
	cfg.ParseENV()
	return cfg.DBName
}

func (cfg *Config) DBPathENV() string {
	cfg.ParseENV()
	return cfg.DBPath
}
