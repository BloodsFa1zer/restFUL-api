package ENV

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	URI_BD string `env:"URI_MongoDB"`
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
	fmt.Println(cfg.URI_BD)
	return cfg.URI_BD
}
