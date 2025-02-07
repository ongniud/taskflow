package ops

import (
	"fmt"

	"github.com/ongniud/taskflow/model"
)

type SubGraphResponse struct {
}

func (op *SubGraphResponse) Execute(ctx model.IOpContext) error {
	fmt.Printf("node[%s.%s] [SubGraphResponse] ctx.GetNodeInputs():%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, ctx.GetInputs())
	fmt.Printf("node[%s.%s] [SubGraphResponse] ctx.GetGraphInputs():%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, ctx.GetGraphInputs())

	for k, v := range map[string]any{
		"3.1": 3.1,
		"3.2": 3.2,
		"3.3": 3.3,
	} {
		ctx.SetGraphOutput(k, v)
	}
	return nil
}
