package main

import (
	"flag"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"tikstart/task"
	"tikstart/task/client/internal/config"
	"time"
)

var configFile = flag.String("f", "etc/client.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	client := asynq.NewClient(
		asynq.RedisClientOpt{Addr: c.Redis.Addr, Password: c.Redis.Pass},
	)

	syncUserCountsTask, err := task.NewSyncUserCountsTask("following_count")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := client.Enqueue(syncUserCountsTask, asynq.MaxRetry(0), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
