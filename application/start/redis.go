package start

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type mRedis struct{
	ctx context.Context
	client *redis.Client
}

func newMRedis(ctx context.Context, client *redis.Client) *mRedis {
	return &mRedis{ctx: ctx, client: client}
}


var MRedisClient *mRedis

func (this *mRedis) GetRedisClient() (client *redis.Client, err error)  {
	if (this.client == nil) {
		return nil, errors.New("redisClient not init")
	}
	return this.client, nil
}

func InitRedis()  {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		errorStr:= "initRedis fail, err:" + err.Error()
		Logger.Error(errorStr)
		panic(errors.New(errorStr))
	}

	MRedisClient = newMRedis(ctx, rdb)
	Logger.Debug("init redis success")
}

func CloseRedis()  {
	if MRedisClient.client != nil {
		MRedisClient.client.Close()
	}
}
