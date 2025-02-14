package taskflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-graphviz"
	"log"
	"time"

	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/tfctx"
	"github.com/ongniud/taskflow/utils"

	"github.com/pkg/errors"
)

// Task is a graph execution unit that runs a graph of nodes with specific execution options, including concurrency control.
// Task also serves as the external interface of the framework, allowing users to interact with it.
type Task struct {
	opts *TaskOptions
	sem  *utils.GoSemaphore
	tc   *tfctx.TaskCtx
}

// NewTask initializes a new Task with the given options.
func NewTask(options ...TaskOption) (*Task, error) {
	opts := &TaskOptions{}
	for _, option := range options {
		option(opts)
	}
	if opts.Graph == nil {
		return nil, fmt.Errorf("nil graph")
	}
	para := opts.Parallelism
	if para == 0 {
		para = DefaultPoolSize
	}
	sem := utils.NewGoSemaphore()
	return &Task{
		opts: opts,
		sem:  sem,
	}, nil
}

// Run executes the task by processing the graph nodes in parallel.
func (e *Task) Run(ctx model.IFlowContext) error {
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	if e == nil {
		return fmt.Errorf("nil task")
	}

	// Create a new task context
	g := e.opts.Graph
	tc, err := tfctx.NewTaskContext(ctx, g)
	if err != nil {
		return err
	}
	e.tc = tc

	// Start execution from the graph's beginning nodes
	tc.NodeInflight(len(g.BeginNodes()))
	for beginNode := range g.BeginNodes() {
		e.sem.Submit(func() {
			e.runNode(tc, tc.GetNodeCtx(beginNode))
			tc.NodeDone()
		})
	}

	// Submit a task to wait for the entire graph to finish execution
	var done = make(chan struct{})
	e.sem.Submit(func() {
		tc.NodeWait()
		close(done)
	})

	// Wait for task completion and handle termination cases
	select {
	case <-tc.Ctx().Done():
		err := tc.Ctx().Err()
		tc.Abort(err)
		tc.NodeWait()
		log.Printf("graph[%s] exit, <-ctx.Done(), err=%v\n", g.Name, err)
		return err
	case <-tc.Aborted():
		err := tc.Err()
		tc.Cancel(err)
		tc.NodeWait()
		log.Printf("graph[%s] exit, execute failed, err=%v\n", g.Name, err)
		return err
	case <-done:
		log.Printf("graph[%s] execute success\n", g.Name)
	}

	ctx.SetOutputs(tc.GetOutputs())
	return nil
}

// runNode executes an individual node in the graph.
func (e *Task) runNode(tc *tfctx.TaskCtx, nc *tfctx.NodeCtx) {
	log.Printf("node[%s.%s] start\n", tc.GetGraph().Name, nc.Node.Name)
	if err := tc.Err(); err != nil {
		log.Printf("node[%s.%s] exit, tc.Error() == %s\n", tc.GetGraph().Name, nc.Node.Name, err.Error())
		nc.Status = tfctx.NodeStatusAbort
		nc.Err = err
		return
	}
	if err := tc.Ctx().Err(); err != nil {
		log.Printf("node[%s.%s] exit, tc.Ctx().Err() == %s\n", tc.GetGraph().Name, nc.Node.Name, err.Error())
		nc.Status = tfctx.NodeStatusAbort
		nc.Err = err
		return
	}

	if nc.Status != tfctx.NodeStatusInit {
		log.Println("nc.Status != tfctx.NodeStatusInit", nc.Status)
		return
	}

	if nc.WaitingPrevCount.Load() != 0 {
		log.Println("nc.WaitingPrevCount.Load() != 0")
		return
	}

	defer func() {
		e.runNextNode(tc, nc)
		if nc.Cancel != nil {
			nc.Cancel()
		}
		log.Printf("node[%s.%s] end\n", tc.GetGraph().Name, nc.Node.Name)
	}()

	if nc.Node.GetInDegree() != 0 && nc.PrunedPrevCount.Load() == int32(nc.Node.GetInDegree()) {
		nc.Status = tfctx.NodeStatusPruned
		log.Printf("node[%s.%s] is pruned\n", tc.GetGraph().Name, nc.Node.Name)
		return
	}

	nc.Status = tfctx.NodeStatusPrepare
	log.Printf("node[%s.%s] is preparing\n", tc.GetGraph().Name, nc.Node.Name)

	inputs, err := e.getNodeInputs(tc, nc)
	if err != nil {
		nc.Status = tfctx.NodeStatusFail
		nc.Err = err
		log.Printf("node[%s.%s] get depend inputs failed\n", tc.GetGraph().Name, nc.Node.Name)
		return
	}

	nc.Start = time.Now()
	nc.Inputs = inputs
	nc.Status = tfctx.NodeStatusRunning
	log.Printf("node[%s.%s] is running\n", tc.GetGraph().Name, nc.Node.Name)

	if nc.Node.Timeout != 0 {
		nc.Ctx, nc.Cancel = context.WithTimeout(tc.Ctx(), time.Duration(nc.Node.Timeout)*time.Millisecond)
	} else {
		nc.Ctx, nc.Cancel = context.WithCancel(tc.Ctx())
	}

	// Execute the node operation or subgraph
	switch nc.Node.Kind {
	case config.NodeKindOperator:
		if err := e.execOperator(tc, nc); err != nil {
			msg := fmt.Sprintf("node[%s.%s] execute operator[%s] failed, err=%s", tc.GetGraph().Name, nc.Node.Name, nc.Node.Operator, err)
			nc.Status = tfctx.NodeStatusFail
			nc.Err = errors.Wrap(err, msg)
			log.Printf("%s\n", msg)
			return
		}
	case config.NodeKindGraph:
		if err := e.execGraph(tc, nc); err != nil {
			msg := fmt.Sprintf("node[%s.%s] execute graph[%s] failed", tc.GetGraph().Name, nc.Node.Name, nc.Node.Graph)
			nc.Status = tfctx.NodeStatusFail
			nc.Err = errors.Wrap(err, msg)
			log.Printf("%s\n", msg)
			return
		}
	default:
		nc.Status = tfctx.NodeStatusFail
		nc.Err = errors.New("unknown node type")
		return
	}

	// Mark node success
	nc.Status = tfctx.NodeStatusSuccess
	nc.Err = nil
	nc.Cost = time.Since(nc.Start)
	log.Printf("node[%s.%s] execute success\n", tc.GetGraph().Name, nc.Node.Name)
	log.Printf("node[%s.%s] stats: %s\n", tc.GetGraph().Name, nc.Node.Name, nc.String())
	return
}

// runNextNode processes the next set of nodes after the current one completes.
func (e *Task) runNextNode(ec *tfctx.TaskCtx, nc *tfctx.NodeCtx) {
	prune := false
	switch nc.Status {
	case tfctx.NodeStatusPruned:
		prune = true
	case tfctx.NodeStatusFail:
		if nc.Node.ErrAbort {
			ec.Abort(nc.Err)
			log.Printf("node[%s.%s] failed abort\n", ec.GetGraph().Name, nc.Node.Name)
			return
		}
		if nc.Node.ErrPrune {
			log.Printf("node[%s.%s] failed prune\n", ec.GetGraph().Name, nc.Node.Name)
			prune = true
		} else {
			log.Printf("node[%s.%s] failed ignore\n", ec.GetGraph().Name, nc.Node.Name)
		}
	}

	readyNodes := make([]*tfctx.NodeCtx, 0)
	for _, next := range nc.Node.GetNextNodes() {
		nextCtx := ec.GetNodeCtx(next)
		if prune {
			nextCtx.PrunedPrevCount.Add(1)
		}
		if nextCtx.WaitingPrevCount.Add(-1) == 0 {
			readyNodes = append(readyNodes, nextCtx)
		}
	}

	ec.NodeInflight(len(readyNodes))
	for _, readyNode := range readyNodes {
		readyNode := readyNode
		log.Printf("node[%s.%s] ready next node: %s\n", ec.GetGraph().Name, nc.Node.Name, readyNode.Node.Name)
		e.sem.Submit(func() {
			e.runNode(ec, readyNode)
			ec.NodeDone()
		})
	}
}

func (e *Task) execOperator(tc *tfctx.TaskCtx, nc *tfctx.NodeCtx) error {
	return nc.Op.Execute(tfctx.NewOpCtx(tc, nc))
}

func (e *Task) execGraph(ec *tfctx.TaskCtx, nc *tfctx.NodeCtx) error {
	inputs := make(map[string]any, len(nc.GetInputs()))
	for _, input := range nc.GetInputs() {
		if input != nil && input.Data != nil {
			field := input.Field.Mapping
			if field == "" {
				field = input.Field.Name
			}
			inputs[field] = input.Data.Val
		}
	}

	if e.opts.Debug {
		inputsMapStr, _ := json.Marshal(inputs)
		log.Printf("node[%s.%s] subgraph begin, inputs=%s\n", ec.GetGraph().Name, nc.Node.Name, string(inputsMapStr))
	}

	task, err := NewTask(WithGraph(nc.SubGrh))
	if err != nil {
		return errors.Wrap(err, "new engine failed")
	}

	flow := tfctx.NewFlowCtx(nc.Ctx).WithInputs(inputs).WithParams(ec.GetParams())
	if err := task.Run(flow); err != nil {
		log.Printf("node[%s.%s] graph execute fail, err=%s\n", ec.GetGraph().Name, nc.Node.Name, err.Error())
		return errors.Wrap(err, "graph execute failed")
	}

	if e.opts.Debug {
		outputsMapStr, _ := json.Marshal(flow.GetOutputs())
		log.Printf("node[%s.%s] subgraph end, outputs=%s\n", ec.GetGraph().Name, nc.Node.Name, string(outputsMapStr))
	}

	ec.SetParams(flow.GetParams())
	nc.SetOutputs(flow.GetOutputs())
	nc.SetTaskCtx(task.tc)
	return nil
}

func (e *Task) getNodeInputs(ec *tfctx.TaskCtx, nc *tfctx.NodeCtx) ([]*model.FieldData, error) {
	refs := nc.Node.GetFieldRefs()
	inputs := make([]*model.FieldData, len(nc.Node.Inputs))
	for i, input := range nc.Node.Inputs {
		ref := refs[i]
		if ref == nil || ref.Idx < 0 {
			log.Printf("node[%s.%s] ref is nil, i=%d, input=%s\n", ec.GetGraph().Name, nc.Node.Name, i, input.Name)
			continue
		}
		node, idx := ref.Node, ref.Idx
		depNodeCtx := ec.GetNodeCtx(node)
		if depNodeCtx == nil {
			log.Printf("node[%s.%s] depend nc nil, i=%d, input=%s, node=%s\n", ec.GetGraph().Name, nc.Node.Name, i, input.Name, node)
			continue
		}

		outputs := depNodeCtx.GetOutputs()
		if idx >= len(outputs) {
			log.Printf("node[%s.%s] input lost due to depend node fail, idx=%d, name=%s, depend node[%s] idx:%d\n", ec.GetGraph().Name, nc.Node.Name, i, input.Name, node, idx)
			continue
		}

		output := outputs[idx]
		if output == nil {
			log.Printf("node[%s.%s] input lost: idx=%d, name=%s, depend node[%s] idx:%d\n", ec.GetGraph().Name, nc.Node.Name, i, input.Name, node, idx)
			if input.Require {
				return nil, errors.New("no required input")
			}
			continue
		}

		log.Printf("node[%s.%s] input found: idx=%d, name=%s, depend node[%s] idx:%d，val=%v\n", ec.GetGraph().Name, nc.Node.Name, i, input.Name, node, idx, output.Data)
		inputs[i] = &model.FieldData{
			Field: input,
			Data:  output.Data,
		}
	}
	return inputs, nil
}

func (e *Task) Ctx() *tfctx.TaskCtx {
	return e.tc
}

func (e *Task) Render() {
	g, _ := graphviz.New(context.Background())
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	vis := NewVisualizer()
	if err := vis.Viz(e.tc, graph); err != nil {
		log.Fatal(err)
	}
	if err := g.RenderFilename(context.Background(), graph, graphviz.PNG, "output.png"); err != nil {
		log.Fatal(err)
	}
}
