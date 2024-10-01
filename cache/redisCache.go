package cache

import (
	"github.com/redis/go-redis/v9"
	"lightRoom/utils"
	"log"
)

var LRedis *redis.Client

func RedisInit() {
	urlRedis := utils.Settings.RedisDsn

	//if urlRedis == "" {
	//	log.Fatal("REDIS_URL env variable not set")
	//}
	redisOptions, err := redis.ParseURL(urlRedis)

	if err != nil {
		log.Fatal(err)
	}
	LRedis = redis.NewClient(redisOptions)
}
