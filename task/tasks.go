package task

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
)

var (
	SyncUserCounts = "user:counts:sync"
)

type SyncPayload struct {
	Field string
}

func NewSync(name string, field string) *asynq.Task {
	payload, err := json.Marshal(SyncPayload{Field: field})
	if err != nil {
		log.Fatalf("cannot marshal payload: %v", err)
	}
	return asynq.NewTask(name, payload)
}
