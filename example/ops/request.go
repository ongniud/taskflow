package ops

import (
	"encoding/json"
	"fmt"

	"github.com/ongniud/taskflow/model"
)

type Request struct{}

func (op *Request) Execute(ctx model.IOpContext) error {
	req, ok := ctx.GetGraphInputs()["request"]
	if !ok {
		fmt.Println("3233")
		return fmt.Errorf("request not exist")
	}
	if req == nil {
		fmt.Println("3233")
		return fmt.Errorf("request is nil")
	}

	reqStr, _ := json.Marshal(req)
	fmt.Println("req:", string(reqStr))
	ctx.SetParam("request", req)

	//ctx.SetGraphOutput("a", 1)
	ctx.SetOutputs(map[string]any{
		"1.1": 1.1,
		"1.2": nil,
		"1.3": 1.3,
	})
	return nil
}
