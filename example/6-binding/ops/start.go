package ops

import (
	"github.com/ongniud/taskflow/model"
)

type Start struct {
}

func (op *Start) Execute(ctx model.IOpContext) error {
	ctx.SetOutputs(map[string]any{
		"request_id":  "r123",
		"user_id":     int64(1234567),
		"session_ids": []int64{123, 123, 123, 123, 123},
		"host_ids":    []int64{123, 123, 123, 123, 123},
	})
	return nil
}
