package tfctx

import (
	"context"
)

type FlowCtx struct {
	context.Context
	inputs  map[string]any
	outputs map[string]any
	params  map[string]any
}

func NewFlowCtx(ctx context.Context) *FlowCtx {
	return &FlowCtx{
		Context: ctx,
	}
}

func (o *FlowCtx) WithInputs(inputs map[string]any) *FlowCtx {
	o.SetInputs(inputs)
	return o
}

func (o *FlowCtx) WithParams(params map[string]any) *FlowCtx {
	o.SetParams(params)
	return o
}

func (o *FlowCtx) GetInputs() map[string]any {
	return o.inputs
}

func (o *FlowCtx) GetOutputs() map[string]any { return o.outputs }

func (o *FlowCtx) GetParams() map[string]any {
	return o.params
}

func (o *FlowCtx) SetInputs(inputs map[string]any) {
	o.inputs = inputs
}
func (o *FlowCtx) SetParams(params map[string]any) {
	o.params = params
}
func (o *FlowCtx) SetOutputs(outputs map[string]any) {
	o.outputs = outputs
	return
}
