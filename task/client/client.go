package main

import (
	"flag"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/conf"
	"tikstart/task"
	"tikstart/task/client/internal/config"
)

var configFile = flag.String("f", "etc/client.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Pass,
		},
		nil,
	)

	entryId, err := scheduler.Register("@every 10s", task.NewSync(task.SyncUserCounts, "following_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("following_count scheduler: %s\n", entryId)

	entryId, err = scheduler.Register("@every 10s", task.NewSync(task.SyncUserCounts, "follower_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("follower_count scheduler: %s\n", entryId)

	entryId, err = scheduler.Register("@every 10s", task.NewSync(task.SyncUserCounts, "total_favorited"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("total_favorited scheduler: %s\n", entryId)

	entryId, err = scheduler.Register("@every 10s", task.NewSync(task.SyncUserCounts, "favorite_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("favorite_count scheduler: %s\n", entryId)

	err = scheduler.Run()
	if err != nil {
		panic(err)
	}
}
