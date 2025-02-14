package tfctx

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/model/graph"
	"github.com/ongniud/taskflow/registry"
)

type TaskCtx struct {
	// Base
	fc *FlowCtx
	g  *graph.Graph

	// Argument
	inputs  map[string]any
	outputs sync.Map
	params  sync.Map

	// Control
	ctx    context.Context
	cancel context.CancelCauseFunc

	// Status
	status int
	wg     sync.WaitGroup
	err    atomic.Value

	abortOnce sync.Once
	abortCh   chan struct{}

	// Data
	nodes map[string]*NodeCtx

	// Statistics
}

func NewTaskContext(ctx model.IFlowContext, g *graph.Graph) (*TaskCtx, error) {
	tc := &TaskCtx{
		g:       g,
		inputs:  make(map[string]any),
		nodes:   make(map[string]*NodeCtx),
		abortCh: make(chan struct{}),
	}
	if err := tc.init(); err != nil {
		return nil, err
	}
	tc.ctx, tc.cancel = context.WithCancelCause(ctx)
	tc.SetInputs(ctx.GetInputs())
	for k, v := range ctx.GetParams() {
		tc.SetParam(k, v)
	}
	return tc, nil
}

func (r *TaskCtx) Aborted() chan struct{} {
	return r.abortCh
}
func (r *TaskCtx) Abort(err error) {
	r.abortOnce.Do(func() {
		r.err.Store(err)
		close(r.abortCh)
	})
}

func (r *TaskCtx) Ctx() context.Context {
	return r.ctx
}
func (r *TaskCtx) Cancel(err error) {
	r.cancel(err)
}

func (r *TaskCtx) Err() error {
	if err := r.err.Load(); err != nil {
		return err.(error)
	}
	return nil
}

func (r *TaskCtx) NodeInflight(cnt int) {
	r.wg.Add(cnt)
}
func (r *TaskCtx) NodeWait() {
	r.wg.Wait()
}
func (r *TaskCtx) NodeDone() {
	r.wg.Done()
}

func (r *TaskCtx) GetGraph() *graph.Graph {
	return r.g
}
func (r *TaskCtx) GetNodeCtx(name string) *NodeCtx {
	return r.nodes[name]
}

func (r *TaskCtx) GetInputs() map[string]any {
	return r.inputs
}
func (r *TaskCtx) SetInputs(inputs map[string]any) {
	r.inputs = inputs
}

func (r *TaskCtx) SetOutput(key string, val any) {
	r.outputs.Store(key, val)
}
func (r *TaskCtx) GetOutputs() map[string]any {
	outputs := make(map[string]any)
	r.outputs.Range(func(key, value any) bool {
		outputs[key.(string)] = value
		return true
	})
	return outputs
}

func (r *TaskCtx) SetParam(key string, val any) {
	r.params.Store(key, val)
}
func (r *TaskCtx) SetParams(params map[string]any) {
	for k, v := range params {
		r.params.LoadOrStore(k, v)
	}
}
func (r *TaskCtx) GetParam(key string) (any, bool) {
	return r.params.Load(key)
}
func (r *TaskCtx) GetParams() map[string]any {
	params := make(map[string]any)
	r.params.Range(func(key, value any) bool {
		params[key.(string)] = value
		return true
	})
	return params
}

func (r *TaskCtx) init() error {
	for _, node := range r.g.Nodes() {
		switch node.Kind {
		case config.NodeKindOperator:
			op, err := registry.GetOp(node.Operator, node.Node)
			if err != nil {
				return err
			}
			r.nodes[node.Name] = NewNodeCtx(r.g, node, op, nil)
		case config.NodeKindGraph:
			g := registry.GetGraph(node.Graph)
			if g == nil {
				return fmt.Errorf("graph not exist")
			}
			r.nodes[node.Name] = NewNodeCtx(r.g, node, nil, g)
		default:
			return fmt.Errorf("unknow node kind")
		}
	}
	return nil
}

func (r *TaskCtx) String() string {
	//	if r.status != NodeStatusSuccess {
	//		return fmt.Sprintf("Node[%s.%s] Status: %d, Err:%s", c.Graph.Name, c.Node.Name, c.Status, c.Err)
	//	}
	//	var is []string
	//	for _, input := range c.Inputs {
	//		if input == nil {
	//			is = append(is, "nil")
	//			continue
	//		}
	//		field := input.Field.Mapping
	//		if input.Data == nil || input.Data.Val == nil {
	//			is = append(is, fmt.Sprintf("<%s, nil>", field))
	//			continue
	//		}
	//		is = append(is, fmt.Sprintf("<%s, %v>", field, input.Data.Val))
	//	}
	//
	//	var os []string
	//	for _, output := range c.Outputs {
	//		if output == nil {
	//			os = append(os, "nil")
	//			continue
	//		}
	//		field := output.Field.Name
	//		if output.Data == nil || output.Data.Val == nil {
	//			os = append(os, fmt.Sprintf("<%s, nil>", field))
	//			continue
	//		}
	//		os = append(os, fmt.Sprintf("<%s, %v>", field, output.Data.Val))
	//	}
	//
	//	return fmt.Sprintf(`Status: %d, Cost: %s
	//Inputs: %s,
	//Outputs: %s,
	//`, c.Status, c.Cost, strings.Join(is, ","), strings.Join(os, ","))
	return ""
}
