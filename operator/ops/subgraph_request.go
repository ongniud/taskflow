package ops

import (
	"fmt"

	"github.com/ongniud/taskflow/model"
)

type SubGraphRequest struct {
}

func (op *SubGraphRequest) Execute(ctx model.IOpContext) error {
	a, ok := ctx.GetGraphInputs()["1.1"]
	if !ok {
		return fmt.Errorf("a not exist")
	}
	b, ok := ctx.GetGraphInputs()["1.2"]
	if !ok {
		return fmt.Errorf("b not exist")
	}
	c, ok := ctx.GetGraphInputs()["1.3"]
	if !ok {
		return fmt.Errorf("c not exist")
	}
	fmt.Printf("node[%s.%s] [SubGraphRequest] ctx.GetNodeInputs():%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, ctx.GetInputs())
	fmt.Printf("node[%s.%s] [SubGraphRequest] ctx.GetGraphInputs():%v\n", ctx.GetGraph().Name, ctx.GetNode().Name, ctx.GetGraphInputs())
	fmt.Printf("node[%s.%s] [SubGraphRequest] abc:%f, %f, %f\n", ctx.GetGraph().Name, ctx.GetNode().Name, a, b, c)
	ctx.SetOutputs(map[string]any{
		"1.1": 1.1,
		"1.2": 1.2,
		"1.3": 1.3,
	})
	return nil
}
