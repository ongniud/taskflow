package ops

import (
	"fmt"
	"github.com/ongniud/taskflow/model"
)

type Process struct{}

func (op *Process) Execute(ctx model.IOpContext) error {
	inputs := ctx.GetInputs()
	ReqId := inputs["req.ReqId"]
	Scene := inputs["req.Scene"]
	Country := inputs["req.Country"]
	DeviceID := inputs["req.DeviceID"]
	Uid := inputs["req.Uid"]
	ctx.SetOutputs(map[string]any{
		"identity": fmt.Sprintf("ReqId:%s,Scene:%s,Country:%s,DeviceId:%s,Uid:%d", ReqId, Scene, Country, DeviceID, Uid),
	})
	return nil
}
