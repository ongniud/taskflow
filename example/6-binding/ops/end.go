package ops

import (
	"github.com/ongniud/taskflow/model"
)

type End struct{}

func (op *End) Execute(ctx model.IOpContext) error {
	inputs := ctx.GetInputs()
	identity := inputs["identity"]
	ctx.SetGraphOutput("identity", identity)
	return nil
}
