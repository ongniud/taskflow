package ops

import (
	"fmt"

	"github.com/ongniud/taskflow/model"
)

type Binding struct {
	Args `tf:"args"`
}

type Args struct {
	ReqId      *string `tf:"input,request_id,required"`
	UserId     *int64  `tf:"input,user_id,optional"`
	SessionIds []int64 `tf:"input,session_ids,_"`
	HostIds    []int64 `tf:"input,host_ids,_"`
	Identity   *string `tf:"output,identity"`
}

func (op *Binding) Execute(ctx model.IOpContext) error {
	fmt.Printf("node[%s.%s] [Debug] args=%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, op.Args)
	args := op.Args
	identity := fmt.Sprintf("ReqId:%s,UserId:%d,SessionIds:%v,HostIds:%v", *args.ReqId, *args.UserId, args.SessionIds, args.HostIds)
	op.Args.Identity = &identity
	return nil
}
