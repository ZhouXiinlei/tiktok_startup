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
	userInterval := GetInterval(c.Interval.UserCounts)
	videoInterval := GetInterval(c.Interval.VideoCounts)
	fmt.Printf("running userCounts sync %s, videoCounts sync %s\n", userInterval, videoInterval)

	// userCounts
	entryId, err := scheduler.Register(userInterval, task.NewSync(task.SyncUserCounts, "following_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("user_following_count scheduler: %s\n", entryId)

	entryId, err = scheduler.Register(userInterval, task.NewSync(task.SyncUserCounts, "follower_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("user_follower_count scheduler: %s\n", entryId)

	entryId, err = scheduler.Register(userInterval, task.NewSync(task.SyncUserCounts, "total_favorited"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("user_total_favorited scheduler: %s\n", entryId)

	entryId, err = scheduler.Register(userInterval, task.NewSync(task.SyncUserCounts, "favorite_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("user_favorite_count scheduler: %s\n", entryId)

	// videoCounts
	entryId, err = scheduler.Register(videoInterval, task.NewSync(task.SyncVideoCounts, "favorite_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("video_favorite_count scheduler: %s\n", entryId)

	entryId, err = scheduler.Register(videoInterval, task.NewSync(task.SyncVideoCounts, "comment_count"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("video_comment_count scheduler: %s\n", entryId)

	// run scheduler
	err = scheduler.Run()
	if err != nil {
		panic(err)
	}
}

func GetInterval(interval string) string {
	return fmt.Sprintf("@every %s", interval)
}
