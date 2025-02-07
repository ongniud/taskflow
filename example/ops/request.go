package ops

import (
	"encoding/json"
	"fmt"
	"github.com/ongniud/taskflow/example/entity"

	"github.com/ongniud/taskflow/model"
)

type Request struct{}

func (op *Request) Execute(ctx model.IOpContext) error {
	request, ok := ctx.GetGraphInputs()["request"]
	if !ok {
		return fmt.Errorf("request not exist")
	}
	if request == nil {
		return fmt.Errorf("request is nil")
	}

	req, ok := request.(*entity.Request)
	if !ok {
		return fmt.Errorf("request is invalid")
	}

	reqStr, _ := json.Marshal(req)
	fmt.Println("req:", string(reqStr))
	ctx.SetParam("request", req)

	ctx.SetOutputs(map[string]any{
		"req.ReqId":    req.ReqId,
		"req.Scene":    req.Scene,
		"req.Country":  req.Country,
		"req.DeviceID": req.DeviceID,
		"req.Uid":      req.Uid,
	})
	return nil
}
