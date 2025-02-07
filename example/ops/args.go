package ops

import (
	"fmt"

	"github.com/ongniud/taskflow/model"
)

type ArgsOp struct {
	Args `tf:"args"`
}

type Args struct {
	ReqId      *string `tf:"input,request_id,required"`
	UserId     *int64  `tf:"input,user_id,optional"`
	SessionIds []int64 `tf:"input,session_ids,_"`
	HostIds    []int64 `tf:"input"`
	Req        *string `tf:"output,req"`
	Ctx        *string `tf:"output,ctx"`
}

func (op *ArgsOp) Execute(ctx model.IOpContext) error {
	fmt.Printf("node[%s.%s] [Debug] args=%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, op.Args)
	return nil
}
