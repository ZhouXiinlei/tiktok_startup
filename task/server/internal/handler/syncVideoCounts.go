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

func (l *TaskHandler) SyncVideoCountsHandler(ctx context.Context, t *asynq.Task) error {
	var payload task.SyncPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}
	fmt.Printf("running sync for field: %s\n", payload.Field)

	zRes, err := l.svcCtx.RDS.ZRangeWithScores(ctx, cache.GenVideoCountsKey(payload.Field), 0, -1).Result()
	if err != nil {
		logx.Error(err.Error())
		return err
	}

	for _, res := range zRes {
		memberStr, ok := res.Member.(string)
		if !ok {
			logx.Error("member type error")
			return err
		}
		member, err := strconv.ParseInt(memberStr, 10, 64)
		if err != nil {
			logx.Error(err.Error())
			return err
		}

		score := int64(res.Score)
		fmt.Printf("topic: %s, member: %d, score: %d\n", payload.Field, member, score)

		err = l.svcCtx.DB.
			Model(&model.Video{}).
			Where("video_id = ?", member).
			UpdateColumn(payload.Field, score).
			Error
		if err != nil {
			logx.Error(err.Error())
			return err
		}

		hit, err := l.svcCtx.RDS.Exists(ctx, cache.GenVideoHeatKey(member)).Result()
		if err != nil {
			logx.Error(err.Error())
			return err
		}
		if hit != 1 {
			fmt.Printf("member %d removed from topic %s\n", member, payload.Field)
			_, err := l.svcCtx.RDS.ZRem(ctx, cache.GenVideoCountsKey(payload.Field), member).Result()
			if err != nil {
				logx.Error(err.Error())
				return err
			}
		}
	}
	return nil
}
