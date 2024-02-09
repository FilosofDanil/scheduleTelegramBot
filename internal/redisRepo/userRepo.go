package redisRepo

import (
	"context"
	"github.com/redis/go-redis/v9"
	"schedulerTelegramBot/configs"
	"strconv"
)

type Session struct {
	ChatId int64  `redis:"chatId"`
	State  string `redis:"state"`
}

type RedisDB struct {
	ctx        *context.Context
	client     *redis.Client
	errCounter int
}

func NewRedisDB(ctx *context.Context, conf configs.RedisConfigs) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Hostname + ":" + conf.Port,
		Password: conf.Password, // no password set
		Username: conf.Username,
	})
	return &RedisDB{ctx: ctx, client: rdb, errCounter: 0}
}

func NewSession(chatId int64, state string) *Session {
	return &Session{ChatId: chatId, State: state}
}

func (rdb *RedisDB) StartReading(key int64, session *Session) {
	var c = *rdb.ctx
	_, err := rdb.client.HSet(c, "session:"+strconv.FormatInt(key, 10), session).Result()

	if err != nil {
		rdb.errCounter++
		if rdb.errCounter >= 5 {
			panic(err)
		}
		go rdb.StartReading(key, session)
		return
	} else {
		rdb.errCounter = 0
	}
	return
}

func (rdb *RedisDB) GetSession(key int64) *Session {
	var c = *rdb.ctx
	var session Session
	err := rdb.client.HGetAll(c, "session:"+strconv.FormatInt(key, 10)).Scan(&session)
	if err != nil {
		rdb.errCounter++
		if rdb.errCounter >= 5 {
			panic(err)
		}
		go rdb.GetSession(key)
	} else {
		rdb.errCounter = 0
	}
	return &session
}
