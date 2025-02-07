package tfctx

import (
	"context"

	"github.com/ongniud/taskflow/model/graph"
	"github.com/ongniud/taskflow/utils"
)

type OpContext struct {
	context.Context
	tc *TaskCtx
	nc *NodeCtx
}

func NewOpCtx(tc *TaskCtx, nc *NodeCtx) *OpContext {
	return &OpContext{
		Context: nc.Ctx,
		tc:      tc,
		nc:      nc,
	}
}

func (o *OpContext) GetInputs() map[string]any {
	if o == nil || o.nc == nil {
		return nil
	}
	return utils.FieldsToMap(o.nc.GetInputs())
}

func (o *OpContext) SetOutputs(outputs map[string]any) {
	o.nc.SetOutputs(outputs)
	return
}

func (o *OpContext) GetParams() map[string]any {
	if o == nil || o.nc == nil {
		return nil
	}
	return nil
}

func (o *OpContext) SetParam(key string, val any) {
	o.tc.SetParam(key, val)
}

func (o *OpContext) GetParam(key string) (any, bool) {
	return o.tc.GetParam(key)
}

func (o *OpContext) GetNode() *graph.Node {
	return o.nc.GetNode()
}

func (o *OpContext) GetGraph() *graph.Graph {
	return o.tc.GetGraph()
}

func (o *OpContext) GetGraphInputs() map[string]any {
	return o.tc.GetInputs()
}

func (o *OpContext) SetGraphOutput(key string, val any) {
	o.tc.SetOutput(key, val)
}
