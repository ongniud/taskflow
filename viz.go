package taskflow

import (
	"fmt"
	"log"
	"strings"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/ongniud/taskflow/model"
	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/tfctx"
)

// Visualizer 是一个用于可视化任务图的工具类
type Visualizer struct{}

// NewVisualizer 创建一个新的 Visualizer 实例
func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

// Viz 可视化任务图
func (v *Visualizer) Viz(taskCtx *tfctx.TaskCtx, vzg *cgraph.Graph) error {
	gvzNodes := make(map[string]*cgraph.Node)
	gvzGraphed := make(map[string]*cgraph.Graph)

	// 遍历所有节点
	for name := range taskCtx.GetGraph().Nodes() {
		node, err := v.createNode(taskCtx, vzg, name)
		if err != nil {
			log.Printf("failed to create node %s: %v\n", name, err)
			continue
		}

		// 处理子图节点
		if err := v.handleSubgraph(taskCtx, vzg, gvzGraphed, name, node); err != nil {
			return err
		}

		// 记录节点
		gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name)] = node
	}

	// 创建节点之间的边
	if err := v.createEdges(taskCtx, vzg, gvzNodes, gvzGraphed); err != nil {
		return err
	}

	return nil
}

// createNode 创建并配置节点
func (v *Visualizer) createNode(taskCtx *tfctx.TaskCtx, vzg *cgraph.Graph, name string) (*cgraph.Node, error) {
	node, err := vzg.CreateNodeByName(fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name))
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	// 设置节点样式
	node.SetShape("record")
	node.SetStyle("filled")
	node.SetFillColor("lightgrey")
	node.SetPenWidth(2)
	node.SetFontColor("black")
	node.SetFontSize(12)
	node.SetFontName("Arial")

	// 设置节点颜色
	if _, ok := taskCtx.GetGraph().BeginNodes()[name]; ok {
		node.SetColor("yellow")
	}
	if _, ok := taskCtx.GetGraph().EndNodes()[name]; ok {
		node.SetColor("blue")
	}

	// 设置节点标签
	ndCtx := taskCtx.GetNodeCtx(name)
	label := v.generateNodeLabel(name, ndCtx)
	node.SetLabel(label)

	return node, nil
}

// handleSubgraph 处理子图节点
func (v *Visualizer) handleSubgraph(taskCtx *tfctx.TaskCtx, vzg *cgraph.Graph, gvzGraphed map[string]*cgraph.Graph, name string, node *cgraph.Node) error {
	ndCtx := taskCtx.GetNodeCtx(name)
	if ndCtx.Node.Kind == config.NodeKindGraph {
		subgraph, err := vzg.SubGraphByName("cluster_" + ndCtx.Node.Graph)
		if err != nil {
			return fmt.Errorf("failed to create subgraph: %w", err)
		}
		//subgraph.SetClusterRank(cgraph.LocalCluster)
		//subgraph.SetStyle(cgraph.FilledGraphStyle)

		// 递归可视化子图
		if err := v.Viz(ndCtx.GetTaskCtx(), subgraph); err != nil {
			return err
		}
		gvzGraphed[ndCtx.Node.Graph] = subgraph
	}
	return nil
}

// createEdges 创建节点之间的边
func (v *Visualizer) createEdges(taskCtx *tfctx.TaskCtx, vzg *cgraph.Graph, gvzNodes map[string]*cgraph.Node, gvzGraphed map[string]*cgraph.Graph) error {
	for name, node := range taskCtx.GetGraph().Nodes() {
		ndCtx := taskCtx.GetNodeCtx(name)
		ndVz := gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name)]
		edgeLabel := v.generateEdgeLabel(ndCtx)

		// 处理子图节点
		if ndCtx.Node.Kind == config.NodeKindGraph {
			grh := gvzGraphed[ndCtx.Node.Graph]
			if grh != nil {
				v.createSubgraphEdges(vzg, ndVz, grh)
			}
		}

		// 创建普通节点之间的边
		prev := gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name)]
		for _, next := range node.GetNextNodes() {
			nxt := gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, next)]
			if err := v.createEdge(vzg, prev, nxt, edgeLabel); err != nil {
				log.Printf("failed to create edge: %v\n", err)
			}
		}
	}
	return nil
}

// createEdge 创建边
func (v *Visualizer) createEdge(vzg *cgraph.Graph, from, to *cgraph.Node, label string) error {
	edge, err := vzg.CreateEdgeByName("", from, to)
	if err != nil {
		return fmt.Errorf("failed to create edge: %w", err)
	}
	edge.SetLabel(label)
	edge.SetColor("red")
	edge.SetStyle(cgraph.DottedEdgeStyle)
	edge.SetArrowHead(cgraph.ODiamondArrow)
	edge.SetPenWidth(2)
	return nil
}

// generateNodeLabel 生成节点的标签
func (v *Visualizer) generateNodeLabel(name string, ndCtx *tfctx.NodeCtx) string {
	statusStr := fmt.Sprintf("status: %s", tfctx.NodeStatusMapping[ndCtx.Status])
	opstr := v.getOperationString(ndCtx)
	errStr := fmt.Sprintf("err: %s", ndCtx.Err)

	if ndCtx.Status != tfctx.NodeStatusSuccess {
		return fmt.Sprintf("{%s|%s\\l|%s\\l|%s\\l}", name, statusStr, opstr, errStr)
	}

	// 生成输入输出信息
	inputPairs := v.generateFieldPairs(ndCtx.GetInputs())
	outputPairs := v.generateFieldPairs(ndCtx.GetOutputs())
	inputStr := fmt.Sprintf("inputs:%s", strings.Join(inputPairs, ";"))
	outputStr := fmt.Sprintf("outputs:%s", strings.Join(outputPairs, ";"))

	return fmt.Sprintf("{%s|%s\\l|%s\\l|%s\\l|%s\\l}", name, statusStr, opstr, inputStr, outputStr)
}

// getOperationString 获取操作字符串
func (v *Visualizer) getOperationString(ndCtx *tfctx.NodeCtx) string {
	switch ndCtx.Node.Kind {
	case config.NodeKindOperator:
		return fmt.Sprintf("op: %s", ndCtx.Node.Operator)
	case config.NodeKindGraph:
		return fmt.Sprintf("graph: %s", ndCtx.Node.Graph)
	default:
		return ""
	}
}

// generateFieldPairs 生成字段键值对
func (v *Visualizer) generateFieldPairs(units []*model.FieldData) []string {
	var pairs []string
	for _, unit := range units {
		if unit == nil || unit.Field == nil {
			continue
		}
		var val any
		if unit.Data != nil {
			val = unit.Data.Val
		}
		pairs = append(pairs, fmt.Sprintf("%s=%v", unit.Field.Name, val))
	}
	return pairs
}

// generateEdgeLabel 生成边的标签
func (v *Visualizer) generateEdgeLabel(ndCtx *tfctx.NodeCtx) string {
	if ndCtx.Status == tfctx.NodeStatusPruned {
		return "prune"
	} else if ndCtx.Status == tfctx.NodeStatusFail {
		if ndCtx.Node.ErrAbort {
			return "abort"
		} else if ndCtx.Node.ErrPrune {
			return "prune"
		} else {
			return "ignore"
		}
	}
	return ""
}

// createSubgraphEdges 创建子图节点之间的边
func (v *Visualizer) createSubgraphEdges(vzg *cgraph.Graph, ndVz *cgraph.Node, grh *cgraph.Graph) {
	fstNd, _ := grh.FirstNode()

	edge1, err := vzg.CreateEdgeByName("", ndVz, fstNd)
	if err != nil {
		log.Println("graph create edge error", err)
	}
	edge1.SetStyle(cgraph.DottedEdgeStyle)
	edge1.SetArrowHead(cgraph.ODiamondArrow)

	edge2, err := vzg.CreateEdgeByName("", fstNd, ndVz)
	if err != nil {
		log.Println("graph create edge error", err)
	}
	edge2.SetStyle(cgraph.DottedEdgeStyle)
	edge2.SetArrowHead(cgraph.EDiamondArrow)
}
