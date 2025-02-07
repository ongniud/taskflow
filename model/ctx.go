package model

import (
	"context"

	"github.com/ongniud/taskflow/model/graph"
)

type IFlowContext interface {
	context.Context
	GetInputs() map[string]any
	SetOutputs(map[string]any)
	GetParams() map[string]any
}

type IOpContext interface {
	context.Context
	GetInputs() map[string]any
	SetOutputs(outputs map[string]any)

	SetParam(key string, val any)
	GetParam(key string) (any, bool)

	GetNode() *graph.Node
	GetGraph() *graph.Graph
	GetGraphInputs() map[string]any
	SetGraphOutput(key string, val any)
}
