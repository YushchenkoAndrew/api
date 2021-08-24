package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func (o *Redis) Init(env *Env) {
	o.Client = redis.NewClient(&redis.Options{
		Addr:     env.RedisHost + ":" + env.RedisPort,
		Password: env.RedisPass,
	})
}

func (o *Redis) Test() {
	ctx := context.Background()
	fmt.Println("Test Redis")
	err := o.Client.Set(ctx, "test", "HELLO WORLD", 0).Err()
	if err != nil {
		panic(err)
	}

	val := o.Client.Get(ctx, "test").String()
	fmt.Println("Get 'test' = ", val)
}
