package redisRepo

import (
	"context"
	"github.com/redis/go-redis/v9"
	"schedulerTelegramBot/configs"
)

type Session struct {
	ChatId int64  `redis:"chatId"`
	State  string `redis:"state"`
}

type RedisDB struct {
	ctx        *context.Context
	client     *redis.Client
	channel    *chan string
	errCounter int
}

func NewRedisDB(ctx *context.Context, ch *chan string, conf configs.RedisConfigs) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Hostname + ":" + conf.Port,
		Password: conf.Password, // no password set
		Username: conf.Username,
	})
	return &RedisDB{ctx: ctx, client: rdb, channel: ch, errCounter: 0}
}

func NewSession(chatId int64, state string) *Session {
	return &Session{ChatId: chatId, State: state}
}

func (rdb *RedisDB) StartReading() {
	var c = *rdb.ctx
	_, err := rdb.client.HSet(c, "session:789", Session{51678, "tgyhuj"}).Result()

	if err != nil {
		rdb.errCounter++
		if rdb.errCounter >= 5 {
			panic(err)
		}
		go rdb.StartReading()
		return
	} else {
		rdb.errCounter = 0
	}
	return
}
