package ops

import (
	"fmt"

	"github.com/ongniud/taskflow/example/5-rpc/entity"
	"github.com/ongniud/taskflow/model"
)

type Request struct{}

func (op *Request) Execute(ctx model.IOpContext) error {
	request, ok := ctx.GetGraphInputs()["request"]
	if !ok {
		return fmt.Errorf("request not exist")
	}
	req, ok := request.(*entity.Request)
	if !ok {
		return fmt.Errorf("request is invalid")
	}
	ctx.SetOutputs(map[string]any{
		"req.ReqId":    req.ReqId,
		"req.Scene":    req.Scene,
		"req.Country":  req.Country,
		"req.DeviceID": req.DeviceID,
		"req.Uid":      req.Uid,
	})
	return nil
}
