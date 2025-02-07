package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ongniud/taskflow"
	"github.com/ongniud/taskflow/example/5-rpc/entity"
	"github.com/ongniud/taskflow/example/5-rpc/ops"
	"github.com/ongniud/taskflow/example/utils"
	"github.com/ongniud/taskflow/tfctx"
)

func main() {
	if err := ops.Init(); err != nil {
		log.Fatal(err)
	}

	g, err := utils.LoadGraph("/Users/jianweili/go/src/github.com/ongniud/taskflow/example/5-rpc/graph.json")
	if err != nil {
		log.Fatalf("读取文件内容失败: %v", err)
	}

	task, err := taskflow.NewTask(taskflow.WithGraph(g))
	if err != nil {
		log.Fatal(err)
	}

	flow := tfctx.NewFlowCtx(context.Background()).
		WithInputs(map[string]any{
			"request": &entity.Request{
				ReqId:    "123456789",
				Scene:    "landing",
				Country:  "America",
				Uid:      101,
				DeviceID: "KE-AS-D1A23-BBD",
			},
		})
	if err := task.Run(flow); err != nil {
		log.Fatal(err)
	}

	outputs := flow.GetOutputs()

	response := outputs["response"]
	responseStr, _ := json.Marshal(response)
	fmt.Println("response:", string(responseStr))
	//task.tc.Render()
}
