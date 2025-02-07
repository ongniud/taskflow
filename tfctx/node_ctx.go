package tfctx

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/graph"
)

const (
	NodeStatusInit    = 0
	NodeStatusPrepare = 1
	NodeStatusRunning = 2
	NodeStatusFail    = 3
	NodeStatusAbort   = 4
	NodeStatusPruned  = 5
	NodeStatusSuccess = 6
)

var NodeStatusMapping = map[int32]string{
	0: "Init",
	1: "Prepare",
	2: "Running",
	3: "Fail",
	4: "Abort",
	5: "Pruned",
	6: "Success",
}

type NodeCtx struct {
	// Base
	Graph  *graph.Graph
	Node   *graph.Node
	Op     model.IOperator
	SubGrh *graph.Graph

	// Runtime
	tc *TaskCtx

	// Status
	Status           int32
	Err              error
	Start            time.Time
	Cost             time.Duration
	ManualOutput     bool
	WaitingPrevCount *atomic.Int32 // atomic
	PrunedPrevCount  *atomic.Int32

	// Control
	Ctx    context.Context
	Cancel context.CancelFunc

	// Data
	Inputs  []*model.FieldData
	Outputs []*model.FieldData

	// Statistics

}

func NewNodeCtx(grh *graph.Graph, node *graph.Node, op model.IOperator, subGrh *graph.Graph) *NodeCtx {
	ndCtx := &NodeCtx{
		Graph:            grh,
		Node:             node,
		Op:               op,
		SubGrh:           subGrh,
		WaitingPrevCount: &atomic.Int32{},
		PrunedPrevCount:  &atomic.Int32{},
	}
	ndCtx.WaitingPrevCount.Store(int32(node.GetInDegree()))
	ndCtx.PrunedPrevCount.Store(0)
	return ndCtx
}

func (c *NodeCtx) GetInputs() []*model.FieldData {
	return c.Inputs
}

func (c *NodeCtx) GetOutputs() []*model.FieldData {
	return c.Outputs
}

func (c *NodeCtx) SetTaskCtx(fc *TaskCtx) {
	c.tc = fc
}

func (c *NodeCtx) GetTaskCtx() *TaskCtx {
	return c.tc
}

func (c *NodeCtx) SetOutputs(outputs map[string]any) {
	fields := make([]*model.FieldData, len(c.Node.Outputs))
	for i, output := range c.Node.Outputs {
		field := output.Mapping
		if field == "" {
			field = output.Name
		}

		val, ok := outputs[field]
		if !ok {
			log.Printf("node[%s.%s] output field lost: idx=%d, name=%s\n", c.Graph.Name, c.Node.Name, i, output.Name)
			continue
		}
		fields[i] = &model.FieldData{
			Field: output,
			Data: &model.Data{
				Val: val,
			},
		}
		log.Printf("node[%s.%s] output field found: idx=%d, name=%s, value=%v\n", c.Graph.Name, c.Node.Name, i, output.Name, val)
	}
	c.Outputs = fields
	c.ManualOutput = true
}

func (c *NodeCtx) GetNode() *graph.Node {
	return c.Node
}

func (c *NodeCtx) String() string {
	if c.Status != NodeStatusSuccess {
		return fmt.Sprintf("Node[%s.%s] Status: %d, Err:%s", c.Graph.Name, c.Node.Name, c.Status, c.Err)
	}
	var is []string
	for _, input := range c.Inputs {
		if input == nil {
			is = append(is, "nil")
			continue
		}
		field := input.Field.Mapping
		if field == "" {
			field = input.Field.Name
		}
		if input.Data == nil || input.Data.Val == nil {
			is = append(is, fmt.Sprintf("<%s, nil>", field))
			continue
		}
		is = append(is, fmt.Sprintf("<%s, %v>", field, input.Data.Val))
	}

	var os []string
	for _, output := range c.Outputs {
		if output == nil {
			os = append(os, "nil")
			continue
		}
		field := output.Field.Name
		if output.Data == nil || output.Data.Val == nil {
			os = append(os, fmt.Sprintf("<%s, nil>", field))
			continue
		}
		os = append(os, fmt.Sprintf("<%s, %v>", field, output.Data.Val))
	}

	return fmt.Sprintf(`Status: %s, Cost: %s
Inputs: %s,
Outputs: %s,
`, NodeStatusMapping[c.Status], c.Cost, strings.Join(is, ","), strings.Join(os, ","))
}
