package task

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

var (
	SyncUserCounts = "user:counts:sync"
)

type SyncUserCountsPayload struct {
	Field string
}

func NewSyncUserCountsTask(field string) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncUserCountsPayload{Field: field})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(SyncUserCounts, payload), nil
}
