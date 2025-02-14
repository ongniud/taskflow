package taskflow

//
//import (
//	"fmt"
//	"log"
//	"strings"
//
//	"github.com/goccy/go-graphviz/cgraph"
//	"github.com/ongniud/taskflow/model"
//	"github.com/ongniud/taskflow/model/config"
//	"github.com/ongniud/taskflow/tfctx"
//)
//
//// Visualizer1 是一个用于可视化任务图的工具类
//type Visualizer1 struct{}
//
//// NewVisualizer1 创建一个新的 Visualizer 实例
//func NewVisualizer1() *Visualizer1 {
//	return &Visualizer1{}
//}
//
//// Viz 可视化任务图
//func (v *Visualizer1) Viz(taskCtx *tfctx.TaskCtx, vzg *cgraph.Graph) error {
//	gvzNodes := make(map[string]*cgraph.Node)
//	gvzGraphed := make(map[string]*cgraph.Graph)
//
//	// 遍历所有节点
//	for name := range taskCtx.GetGraph().Nodes() {
//		nd, err := vzg.CreateNode(fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name))
//		if err != nil {
//			log.Println("graph create node error")
//			continue
//		}
//		nd.SetShape("record")
//
//		// 设置节点颜色
//		if _, ok := taskCtx.GetGraph().BeginNodes()[name]; ok {
//			nd.SetColor("yellow")
//		}
//		if _, ok := taskCtx.GetGraph().EndNodes()[name]; ok {
//			nd.SetColor("blue")
//		}
//
//		// 获取节点上下文
//		ndCtx := taskCtx.GetNodeCtx(name)
//		var opstr string
//		switch ndCtx.Node.Kind {
//		case config.NodeKindOperator:
//			opstr = fmt.Sprintf("op: %s", ndCtx.Node.Operator)
//		case config.NodeKindGraph:
//			opstr = fmt.Sprintf("graph: %s", ndCtx.Node.Graph)
//			subgraph := vzg.SubGraph("cluster_"+ndCtx.Node.Graph, 1)
//			if subgraph == nil {
//				fmt.Println("!!!!")
//			}
//			subgraph.SetClusterRank(cgraph.LocalCluster)
//			subgraph.SetStyle(cgraph.DottedGraphStyle)
//
//			// 递归可视化子图
//			if err := v.Viz(ndCtx.GetTaskCtx(), subgraph); err != nil {
//				return err
//			}
//			gvzGraphed[ndCtx.Node.Graph] = subgraph
//		}
//
//		// 设置节点标签
//		statusStr := fmt.Sprintf("status: %s", tfctx.NodeStatusMapping[ndCtx.Status])
//		errStr := fmt.Sprintf("err: %s", ndCtx.Err)
//		label := v.generateNodeLabel(name, statusStr, opstr, errStr, ndCtx)
//		nd.SetLabel(label)
//
//		// 记录节点
//		gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name)] = nd
//	}
//
//	// 创建节点之间的边
//	for name, node := range taskCtx.GetGraph().Nodes() {
//		ndCtx := taskCtx.GetNodeCtx(name)
//		ndVz := gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name)]
//		edgeLabel := v.generateEdgeLabel(ndCtx)
//
//		// 处理子图节点
//		if ndCtx.Node.Kind == config.NodeKindGraph {
//			grh := gvzGraphed[ndCtx.Node.Graph]
//			if grh != nil {
//				v.createSubgraphEdges(vzg, ndVz, grh)
//			}
//		}
//
//		// 创建普通节点之间的边
//		prev := gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, name)]
//		for _, next := range node.GetNextNodes() {
//			nxt := gvzNodes[fmt.Sprintf("%s.%s", taskCtx.GetGraph().Name, next)]
//			edge, err := vzg.CreateEdge("", prev, nxt)
//			if err != nil {
//				log.Println("graph create edge error", err)
//			}
//			edge.SetLabel(edgeLabel)
//		}
//	}
//	return nil
//}
//
//// generateNodeLabel 生成节点的标签
//func (v *Visualizer1) generateNodeLabel(name, statusStr, opstr, errStr string, ndCtx *tfctx.NodeCtx) string {
//	if ndCtx.Status != tfctx.NodeStatusSuccess {
//		return fmt.Sprintf("{%s|%s\\l|%s\\l|%s\\l}", name, statusStr, opstr, errStr)
//	}
//
//	// 生成输入输出信息
//	inputPairs := v.generateFieldPairs(ndCtx.GetInputs())
//	outputPairs := v.generateFieldPairs(ndCtx.GetOutputs())
//	inputStr := fmt.Sprintf("inputs:%s", strings.Join(inputPairs, ";"))
//	outputStr := fmt.Sprintf("outputs:%s", strings.Join(outputPairs, ";"))
//
//	return fmt.Sprintf("{%s|%s\\l|%s\\l|%s\\l|%s\\l}", name, statusStr, opstr, inputStr, outputStr)
//}
//
//// generateFieldPairs 生成字段键值对
//func (v *Visualizer1) generateFieldPairs(units []*model.FieldData) []string {
//	var pairs []string
//	for _, unit := range units {
//		if unit == nil || unit.Field == nil {
//			continue
//		}
//		var val any
//		if unit.Data != nil {
//			val = unit.Data.Val
//		}
//		pairs = append(pairs, fmt.Sprintf("%s=%v", unit.Field.Name, val))
//	}
//	return pairs
//}
//
//// generateEdgeLabel 生成边的标签
//func (v *Visualizer1) generateEdgeLabel(ndCtx *tfctx.NodeCtx) string {
//	if ndCtx.Status == tfctx.NodeStatusPruned {
//		return "prune"
//	} else if ndCtx.Status == tfctx.NodeStatusFail {
//		if ndCtx.Node.ErrAbort {
//			return "abort"
//		} else if ndCtx.Node.ErrPrune {
//			return "prune"
//		} else {
//			return "ignore"
//		}
//	}
//	return ""
//}
//
//// createSubgraphEdges 创建子图节点之间的边
//func (v *Visualizer1) createSubgraphEdges(vzg *cgraph.Graph, ndVz *cgraph.Node, grh *cgraph.Graph) {
//	edge1, err := vzg.CreateEdge("", ndVz, grh.FirstNode())
//	if err != nil {
//		log.Println("graph create edge error", err)
//	}
//	edge1.SetStyle(cgraph.DottedEdgeStyle)
//	edge1.SetArrowHead(cgraph.ODiamondArrow)
//
//	edge2, err := vzg.CreateEdge("", grh.FirstNode(), ndVz)
//	if err != nil {
//		log.Println("graph create edge error", err)
//	}
//	edge2.SetStyle(cgraph.DottedEdgeStyle)
//	edge2.SetArrowHead(cgraph.EDiamondArrow)
//}
