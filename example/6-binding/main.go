package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ongniud/taskflow"
	"github.com/ongniud/taskflow/example/6-binding/ops"
	"github.com/ongniud/taskflow/example/utils"
	"github.com/ongniud/taskflow/tfctx"
)

func main() {
	if err := ops.Init(); err != nil {
		log.Fatal(err)
	}

	g, err := utils.LoadGraph("example/6-binding/graph.json")
	if err != nil {
		log.Fatalf("load graph failed: %v", err)
	}

	task, err := taskflow.NewTask(taskflow.WithGraph(g))
	if err != nil {
		log.Fatal(err)
	}

	flow := tfctx.NewFlowCtx(context.Background()).
		WithInputs(map[string]any{
			"uid":     13579,
			"age":     18,
			"country": "china",
		}).
		WithParams(map[string]any{})
	if err := task.Run(flow); err != nil {
		log.Fatal(err)
	}

	outputs := flow.GetOutputs()
	outputsStr, _ := json.Marshal(outputs)
	fmt.Println("outputs:", string(outputsStr))
	task.Render("example/6-binding/graph.dot")
}
