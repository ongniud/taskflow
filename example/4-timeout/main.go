package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ongniud/taskflow/example/ops"
	"io"
	"log"
	"os"

	"github.com/ongniud/taskflow"
	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/model/graph"
	"github.com/ongniud/taskflow/tfctx"
)

func main() {
	if err := ops.Init(); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("/Users/jianweili/go/src/github.com/ongniud/taskflow/example/sequence/graph.json")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件内容失败: %v", err)
	}

	var cfg config.Graph
	if err = json.Unmarshal(content, &cfg); err != nil {
		log.Fatalf("反序列化失败: %v", err)
	}

	gra, err := graph.NewGraph(&cfg)
	if err != nil {
		log.Fatalf("new graph fail: %v", err)
	}

	task, err := taskflow.NewTask(taskflow.WithGraph(gra))
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
