package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DbName     string `env:"DatabaseName"`
	DbPath     string `env:"DatabasePath"`
	SigningKey string `env:"SigningKey"`
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

type JwtCustomClaims struct {
	ID   int64  `json:"ID"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewConfig() echojwt.Config {
	cfg := LoadENV("config/.env")
	cfg.ParseENV()

	Config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(cfg.SigningKey),
	}
	return Config
}
