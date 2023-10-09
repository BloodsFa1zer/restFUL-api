package redisDB

import (
	"app3.1/config"
	"app3.1/database"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

// Maybe no new directory needed? and just put it in the database directory?
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
		log.Warn().Err(err).Msg("can`t connect to redisDB server")
		return nil
	}

	return &ClientRedis{Connection: client}
}

func (cl *ClientRedis) GetUser(ID int64) (*database.User, error) {

	resultJSON, err := cl.Connection.Get(strconv.FormatInt(ID, 10)).Result()
	if err != nil {
		log.Warn().Err(err).Msg("can`t get user from Redis cache")
		return nil, err
	}
	var retrievedUser database.User
	err = json.Unmarshal([]byte(resultJSON), &retrievedUser)
	if err != nil {
		log.Warn().Err(err).Msg("can`t unmarshal data from redis cache")
	}

	return &retrievedUser, nil
}

func (cl *ClientRedis) GetUsers() (*[]database.User, error) {
	resultJSON, err := cl.Connection.Get("users").Result()
	if err != nil {
		log.Warn().Err(err).Msg("can`t get user from Redis cache")
		return nil, err
	}

	var users []database.User
	err = json.Unmarshal([]byte(resultJSON), &users)
	if err != nil {
		log.Warn().Err(err).Msg("can`t unmarshal data from redis cache")
	}

	return &users, nil
}

func (cl *ClientRedis) SetUser(user database.User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Warn().Err(err).Msg("can`t marshal data for the further caching")
		return err
	}

	err = cl.Connection.Set(strconv.FormatInt(user.ID, 10), userJSON, 1*time.Minute).Err()
	if err != nil {
		log.Warn().Err(err).Msg("can`t set user for Redis cache")
		return err
	}

	return nil
}

func (cl *ClientRedis) SetUsers(users []database.User) error {
	userJSON, err := json.Marshal(users)
	if err != nil {
		log.Warn().Err(err).Msg("can`t marshal data for the further caching")
		return err
	}

	err = cl.Connection.Set("users", userJSON, 1*time.Minute).Err()
	if err != nil {
		log.Warn().Err(err).Msg("can`t set user for Redis cache")
		return err
	}

	return nil
}
