package task

import (
	"sync"

	"github.com/hibiken/asynq"
)

var (
	client *asynq.Client
	once   sync.Once
)

func Init(redisAddress string) {
	once.Do(func() {
		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr: redisAddress,
			
		})
	})
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func GetClient() *asynq.Client {
	return client
}
