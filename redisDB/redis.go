package redisDB

import (
	"app3.1/config"
	"app3.1/database"
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type ClientRedis struct {
	Connection      *redis.Client
	CacheConnection *cache.Cache
}

func NewClientRedis() *ClientRedis {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	num, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       num,
	})

	_, err = client.Ping(context.TODO()).Result()
	if err != nil {
		log.Warn().Err(err).Msg("can`t connect to redisDB server")
		return nil
	}

	myCache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &ClientRedis{Connection: client, CacheConnection: myCache}
}

func (cl *ClientRedis) GetUser(ID int64) (*database.User, error) {
	var user database.User

	if cl.CacheConnection.Exists(context.TODO(), strconv.FormatInt(ID, 10)) {
		err := cl.CacheConnection.Get(context.TODO(), strconv.FormatInt(ID, 10), &user)
		if err != nil {
			log.Warn().Err(err).Msg("can`t get user from Redis cache")
			return nil, err
		}
	} else {
		return nil, redis.Nil
	}

	return &user, nil
}

func (cl *ClientRedis) GetUsers() (*[]database.User, error) {
	var users []database.User
	if cl.CacheConnection.Exists(context.TODO(), "users") {
		err := cl.CacheConnection.Get(context.TODO(), "users", &users)
		if err != nil {
			log.Warn().Err(err).Msg("can`t get user from Redis cache")
			return nil, err
		}
	} else {
		return nil, redis.Nil
	}
	return &users, nil
}

func (cl *ClientRedis) SetUser(user database.User) error {
	err := cl.CacheConnection.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   strconv.FormatInt(user.ID, 10),
		Value: user,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Warn().Err(err).Msg("can`t set data")
		return err
	}

	return nil
}

func (cl *ClientRedis) SetUsers(users []database.User) error {
	err := cl.CacheConnection.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   "users",
		Value: users,
		TTL:   time.Minute,
	})
	if err != nil {
		return err
	}

	return nil
}
