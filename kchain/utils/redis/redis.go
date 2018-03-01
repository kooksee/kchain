package types

import (
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/tendermint/tmlibs/log"
	"sync"
)

var (
	once sync.Once
	instance *RedisClient
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "redis")
)

type RedisClient struct {
	cfg        *redis.Options
	client     *redis.Client
	isRetrying bool
}

func Init(cfg *redis.Options) {
	once.Do(func() {
		instance = &RedisClient{
			cfg:cfg,
		}
		client := instance.conn()
		pong, err := client.Ping().Result()
		if pong == "" || err != nil {
			panic(err.Error())
		}
	})
}

func GetRedisClient() *redis.Client {
	return instance.client
}

func GetDefault() *RedisClient {
	r := &RedisClient{
		cfg:&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}
	r.conn()
	return r
}

func (r *RedisClient) ParseURL(redisURL string) (*redis.Options, error) {
	return redis.ParseURL(redisURL)
}

func (r *RedisClient) conn() *redis.Client {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.cfg.Addr,
		Password: r.cfg.Password,
		DB:       r.cfg.DB,
	})
	return r.client
}

func (r *RedisClient) ReConn() {
	if r.isRetrying {
		return
	}

	go func() {
		r.isRetrying = true
		for {
			client := r.conn()
			pong, err := client.Ping().Result()
			if pong == "" || err != nil {
				logger.Error(err.Error())
				time.Sleep(2 * time.Second)
			}
			r.isRetrying = false
			return
		}
	}()
}




