package tfctx

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
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

func (r *TaskCtx) Viz(vzg *cgraph.Graph) error {
	gvzNodes := make(map[string]*cgraph.Node)
	gvzGraphed := make(map[string]*cgraph.Graph)
	for name := range r.g.Nodes() {
		nd, err := vzg.CreateNode(fmt.Sprintf("%s.%s", r.GetGraph().Name, name))
		if err != nil {
			log.Println("graph create node error")
			continue
		}
		nd.SetShape("record")
		if _, ok := r.g.BeginNodes()[name]; ok {
			nd.SetColor("yellow")
		}
		if _, ok := r.g.EndNodes()[name]; ok {
			nd.SetColor("blue")
		}

		ndCtx := r.GetNodeCtx(name)
		var opstr string
		switch ndCtx.Node.Kind {
		case config.NodeKindOperator:
			opstr = fmt.Sprintf("op: %s", ndCtx.Node.Operator)
		case config.NodeKindGraph:
			opstr = fmt.Sprintf("graph: %s", ndCtx.Node.Graph)
			subgraph := vzg.SubGraph("cluster_"+ndCtx.Node.Graph, 1)
			if subgraph == nil {
				fmt.Println("!!!!")
			}
			subgraph.SetClusterRank(cgraph.LocalCluster)
			subgraph.SetStyle(cgraph.DottedGraphStyle)

			fmt.Println("!!!!22222")
			if err := ndCtx.GetTaskCtx().Viz(subgraph); err != nil {
				return err
			}
			gvzGraphed[ndCtx.Node.Graph] = subgraph
		}

		statusStr := fmt.Sprintf("status: %s", NodeStatusMapping[ndCtx.Status])
		errStr := fmt.Sprintf("err: %s", ndCtx.Err)

		var label string
		if ndCtx.Status != NodeStatusSuccess {
			label = fmt.Sprintf("{%s|%s\\l|%s\\l|%s\\l}", name, statusStr, opstr, errStr)
		} else {
			var inputPairs []string
			for _, unit := range ndCtx.GetInputs() {
				if unit == nil || unit.Field == nil {
					continue
				}
				var val any
				if unit.Data != nil {
					val = unit.Data.Val
				}
				kv := fmt.Sprintf("%s=%v", unit.Field.Name, val)
				inputPairs = append(inputPairs, kv)
			}
			inputStr := fmt.Sprintf("inputs:%s", strings.Join(inputPairs, ";"))

			var outputPairs []string
			for _, unit := range ndCtx.GetOutputs() {
				if unit == nil || unit.Field == nil {
					continue
				}
				var val any
				if unit.Data != nil {
					val = unit.Data.Val
				}
				kv := fmt.Sprintf("%s=%v", unit.Field.Name, val)
				outputPairs = append(outputPairs, kv)
			}
			outputStr := fmt.Sprintf("outputs:%s", strings.Join(outputPairs, ";"))
			label = fmt.Sprintf("{%s|%s\\l|%s\\l|%s\\l|%s\\l}", name, statusStr, opstr, inputStr, outputStr)
		}
		nd.SetLabel(label)
		gvzNodes[fmt.Sprintf("%s.%s", r.GetGraph().Name, name)] = nd
	}

	for name, node := range r.g.Nodes() {
		ndCtx := r.GetNodeCtx(name)
		ndVz := gvzNodes[fmt.Sprintf("%s.%s", r.GetGraph().Name, name)]
		var edgeLabel string
		if ndCtx.Status == NodeStatusPruned {
			edgeLabel = "prune"
		} else if ndCtx.Status == NodeStatusFail {
			if ndCtx.Node.ErrAbort {
				edgeLabel = "abort"
			} else if ndCtx.Node.ErrPrune {
				edgeLabel = "prune"
			} else {
				edgeLabel = "ignore"
			}
		}

		if ndCtx.Node.Kind == config.NodeKindGraph {
			grh := gvzGraphed[ndCtx.Node.Graph]
			if grh == nil {
				fmt.Println("33333")
			}
			if vzg == nil {
				fmt.Println("44444")
			}

			fmt.Println("55555", grh.FirstNode())
			edge1, err := vzg.CreateEdge("", ndVz, grh.FirstNode())
			if err != nil {
				log.Println("graph create edge error", err)
			}
			edge1.SetStyle(cgraph.DottedEdgeStyle)
			edge1.SetArrowHead(cgraph.ODiamondArrow)

			edge2, err := vzg.CreateEdge("", grh.FirstNode(), ndVz)
			if err != nil {
				log.Println("graph create edge error", err)
			}
			edge2.SetStyle(cgraph.DottedEdgeStyle)
			edge2.SetArrowHead(cgraph.EDiamondArrow)
		}

		fmt.Println(name, ndCtx.Status, ndCtx.Node.ErrAbort, ndCtx.Node.ErrIgnore)
		fmt.Println(edgeLabel)
		prev := gvzNodes[fmt.Sprintf("%s.%s", r.GetGraph().Name, name)]
		for _, next := range node.GetNextNodes() {
			nxt := gvzNodes[fmt.Sprintf("%s.%s", r.GetGraph().Name, next)]
			edge, err := vzg.CreateEdge("", prev, nxt)
			if err != nil {
				log.Println("graph create edge error", err)
			}
			edge.SetLabel(edgeLabel)
		}
	}
	return nil
}

func (r *TaskCtx) Render() {
	gvz := graphviz.New()
	defer gvz.Close()

	vizGraph, _ := gvz.Graph()
	defer vizGraph.Close()
	vizGraph.SetLayout(string(graphviz.DOT))
	vizGraph.SetRankDir("TB")

	if err := r.Viz(vizGraph); err != nil {
		log.Println("failed to create file", err)
	}

	buf := new(bytes.Buffer)
	if err := gvz.Render(vizGraph, graphviz.Format(graphviz.PNG), buf); err != nil {
		log.Println("vizGraph render error", err)
	}

	// 保存图片到文件
	fileName := "graph.png"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("failed to create file", err)
		return
	}
	defer file.Close()

	if _, err := buf.WriteTo(file); err != nil {
		log.Println("failed to write image to file", err)
		return
	}
	// 打开生成的图片文件
	cmd := exec.Command("open", fileName)
	if err := cmd.Start(); err != nil {
		log.Println("failed to open image file", err)
	}
}
