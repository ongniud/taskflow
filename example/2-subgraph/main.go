package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ongniud/taskflow"
	"github.com/ongniud/taskflow/example/2-subgraph/ops"
	"github.com/ongniud/taskflow/example/utils"
	"github.com/ongniud/taskflow/registry"
	"github.com/ongniud/taskflow/tfctx"
)

func main() {
	if err := ops.Init(); err != nil {
		log.Fatal(err)
	}
	subGra1, err := utils.LoadGraph("example/2-subgraph/graph1.json")
	if err != nil {
		log.Fatalf("load graph failed: %v", err)
	}
	subGra2, err := utils.LoadGraph("example/2-subgraph/graph2.json")
	if err != nil {
		log.Fatalf("load graph failed: %v", err)
	}
	if err := registry.RegisterGraph(subGra1); err != nil {
		log.Fatalf("register graph failed: %v", err)
	}
	if err := registry.RegisterGraph(subGra2); err != nil {
		log.Fatalf("register graph failed: %v", err)
	}

	g, err := utils.LoadGraph("example/2-subgraph/graph.json")
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
	task.Render("example/2-subgraph/graph.dot")
}
