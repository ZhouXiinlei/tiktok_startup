package main

import (
	"context"
	"flag"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"tikstart/task"
	"tikstart/task/server/internal/config"
	"tikstart/task/server/internal/handler"
	"tikstart/task/server/internal/svc"
)

var configFile = flag.String("f", "etc/server.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	ctx := context.Background()
	svcCtx := svc.NewServiceContext(c)
	taskHandler := handler.NewTaskHandler(ctx, svcCtx)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Addr, Password: c.Redis.Pass},
		asynq.Config{},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.SyncUserCounts, taskHandler.SyncUserCountsHandler)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
