package ops

import (
	"encoding/json"
	"fmt"

	"github.com/ongniud/taskflow/model"
)

type Debug struct {
	DebugArgs `tf:"args"`
}

type DebugArgs struct {
	//ReqId      *string `tf:"input,request_id,required"`
	//UserId     *int64  `tf:"input,user_id,optional"`
	//SessionIds []int64 `tf:"input,session_ids,_"`
	//HostIds    []int64 `tf:"input"`
	//Req        *string `tf:"output,req"`
	//Ctx        *string `tf:"output,ctx"`
}

func (op *Debug) Execute(ctx model.IOpContext) error {
	param := ctx.GetNode().GetParam()
	var ps map[string]any
	if err := json.Unmarshal([]byte(param), &ps); err != nil {
		return err
	}
	inputs := ctx.GetInputs()
	fmt.Printf("node[%s.%s] [Debug] inputs=%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, inputs)
	ctx.SetOutputs(ps)
	return nil
}
