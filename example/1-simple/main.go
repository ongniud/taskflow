package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ongniud/taskflow"
	"github.com/ongniud/taskflow/example/utils"
	"github.com/ongniud/taskflow/operator"
	"github.com/ongniud/taskflow/tfctx"
)

func main() {
	if err := operator.Init(); err != nil {
		log.Fatal(err)
	}

	g, err := utils.LoadGraph("/Users/jianweili/go/src/github.com/ongniud/taskflow/example/1-simple/graph.json")
	if err != nil {
		log.Fatalf("读取文件内容失败: %v", err)
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
	//task.tc.Render()
}
