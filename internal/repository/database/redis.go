package database

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

func NewRedisDB() (*redis.Client, time.Duration, error) {
	host := viper.GetString("database.redis.host")
	port := viper.GetString("database.redis.port")
	dbname := viper.GetInt("database.redis.name")
	ttl := viper.GetInt("database.redis.ttl")

	RDB := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
		DB:   dbname,
	})

	return RDB, time.Duration(ttl), nil
}
