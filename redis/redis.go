package redis

import (
	"app3.1/config"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type ClientRedis struct {
	Connection *redis.Client
}

func NewClientRedis() *ClientRedis {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Warn().Err(err).Msg("can`t connect to redis server")
		return nil
	}

	return &ClientRedis{Connection: client}
}

//func (cr *ClientRedis) Set() {
//
//	database.User{}
//	err = client.Set("1", database.User{}, 0).Err()
//	// if there has been an error setting the value
//	// handle the error
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	val, err := client.Get("name").Result()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(val)
//}
