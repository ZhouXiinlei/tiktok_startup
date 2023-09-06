package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"tikstart/common/cache"
	"tikstart/common/model"
	"tikstart/task"
)

func (l *TaskHandler) SyncUserCountsHandler(ctx context.Context, t *asynq.Task) error {
	var payload task.SyncPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}
	fmt.Printf("running sync for field: %s\n", payload.Field)

	members, err := l.svcCtx.RDS.ZRange(ctx, cache.GenUserCountsKey(payload.Field), 0, -1).Result()
	if err != nil {
		logx.Error(err.Error())
		return err
	}

	for _, member := range members {
		score, err := l.svcCtx.RDS.ZScore(ctx, cache.GenUserCountsKey(payload.Field), member).Result()
		if err != nil {
			logx.Error(err.Error())
			return err
		}
		fmt.Printf("topic: %s, member: %s, score: %.f\n", payload.Field, member, score)

		err = l.svcCtx.DB.
			Model(&model.User{}).
			Where("user_id = ?", member).
			UpdateColumn(payload.Field, score).
			Error
		if err != nil {
			logx.Error(err.Error())
			return err
		}

		num, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			logx.Error(err.Error())
			return err
		}

		hit, err := l.svcCtx.RDS.Exists(ctx, cache.GenUserHeatKey(num)).Result()
		if err != nil {
			logx.Error(err.Error())
			return err
		}
		if hit != 1 {
			fmt.Printf("member %s removed from topic %s\n", member, payload.Field)
			_, err := l.svcCtx.RDS.ZRem(ctx, cache.GenUserCountsKey(payload.Field), member).Result()
			if err != nil {
				logx.Error(err.Error())
				return err
			}
		}
	}
	return nil
}
