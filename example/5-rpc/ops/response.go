package ops

import (
	"github.com/ongniud/taskflow/example/5-rpc/entity"
	"github.com/ongniud/taskflow/model"
)

type Response struct{}

func (op *Response) Execute(ctx model.IOpContext) error {
	inputs := ctx.GetInputs()
	identity := inputs["identity"]
	rsp := &entity.Response{
		Identity: identity.(string),
	}
	ctx.SetGraphOutput("response", rsp)
	return nil
}
