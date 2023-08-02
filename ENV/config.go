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

func LoadENV(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Panic().Err(err).Msg(" does not load .ENV")
	}
	log.Info().Msg("successfully load .ENV")

}

func ParseENV() *Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panic().Err(err).Msg(" unable to parse environment variables")
	}
	log.Info().Msg("successfully parsed .ENV")
	return &cfg
}

func MongoENV() string {
	cfg := ParseENV()
	fmt.Println(cfg.URI_BD)
	return cfg.URI_BD
}
