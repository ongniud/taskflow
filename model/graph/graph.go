package graph

import (
	"errors"
	"fmt"

	"github.com/ongniud/taskflow/model/config"
)

type Graph struct {
	*config.Graph

	nodes      map[string]*Node       // Map to store nodes by their name
	beginNodes map[string]struct{}    // Set of nodes with no predecessors (start nodes)
	endNodes   map[string]struct{}    // Set of nodes with no successors (end nodes)
	outputs    map[string][]*FieldRef // Map to store output field references for nodes
}

// NewGraph initializes a new Graph from a given config.Graph
func NewGraph(g *config.Graph) (*Graph, error) {
	if g == nil {
		return nil, errors.New("g is nil")
	}
	graph := &Graph{
		Graph:      g,
		nodes:      make(map[string]*Node),
		outputs:    make(map[string][]*FieldRef),
		beginNodes: make(map[string]struct{}),
		endNodes:   make(map[string]struct{}),
	}
	for _, nd := range g.Nodes {
		if err := graph.addNode(NewNode(nd)); err != nil {
			return nil, err
		}
	}
	if err := graph.build(); err != nil {
		return nil, err
	}
	return graph, nil
}

func (g *Graph) Nodes() map[string]*Node {
	return g.nodes
}

func (g *Graph) BeginNodes() map[string]struct{} {
	return g.beginNodes
}

func (g *Graph) EndNodes() map[string]struct{} {
	return g.endNodes
}

// addNode adds a new node to the graph
func (g *Graph) addNode(node *Node) error {
	if node == nil || node.Name == "" {
		return fmt.Errorf("node is nil")
	}
	if _, exists := g.nodes[node.Name]; exists {
		return fmt.Errorf("node[%s] already exists", node.Name)
	}
	if node.Kind == "" {
		node.Kind = config.NodeKindOperator
	}
	// Add node outputs to the outputs map
	for i, output := range node.Outputs {
		field := output.Name
		if _, ok := g.outputs[field]; !ok {
			g.outputs[field] = make([]*FieldRef, 0, 1)
		}
		g.outputs[field] = append(g.outputs[field], &FieldRef{
			Node: node.Name,
			Idx:  i,
		})
	}
	g.nodes[node.Name] = node
	return nil
}

// build constructs the graph edges and checks for loops
func (g *Graph) build() error {
	if err := g.buildEdges(); err != nil {
		return err
	}
	if g.hasLoop() {
		return fmt.Errorf("has loop")
	}
	for _, node := range g.nodes {
		if len(node.GetPrevNodes()) == 0 {
			g.beginNodes[node.Name] = struct{}{}
		}
		if len(node.GetNextNodes()) == 0 {
			g.endNodes[node.Name] = struct{}{}
		}
	}
	return nil
}

// buildEdges creates edges between nodes based on their inputs and outputs
func (g *Graph) buildEdges() error {
	for _, node := range g.nodes {
		for _, field := range node.Inputs {
			if field.Node == "" && field.Name == "" {
				node.AddFieldRef(nil)
				continue
			}

			// If no field name is specified, a node name must be specified
			if field.Name == "" {
				prev, exists := g.nodes[field.Node]
				if !exists {
					return fmt.Errorf("node[%s]'s prev node[%s] not found", node.Name, field.Node)
				}
				prev.AddNext(node.Name) // build edge: prev -> node
				node.AddPrev(prev.Name) // build edge: node -> prev
				node.AddFieldRef(&FieldRef{
					Node: node.Name,
					Idx:  -1,
				})
				continue
			}

			// Find field in global outputs
			refers, ok := g.outputs[field.Name]
			if !ok {
				return fmt.Errorf("node[%s]'s depend field[%s] not found", node.Name, field.Name)
			}

			// If no node name is specified, ensure there is only one candidate
			if field.Node == "" {
				if len(refers) != 1 {
					return fmt.Errorf("node[%s]'s depend field[%s] have more than one candidate[%v]", node.Name, field.Name, refers)
				}
				dependNode := g.nodes[refers[0].Node]
				dependNode.AddNext(node.Name) // build edge: prev -> node
				node.AddPrev(refers[0].Node)  // build edge: node -> prev
				node.AddFieldRef(refers[0])
				continue
			}

			// If node name is specified, find the reference
			var ref *FieldRef
			for _, refer := range refers {
				if refer.Node == field.Node {
					ref = refer
					break
				}
			}

			// Report error if reference is not found
			if ref == nil {
				return fmt.Errorf("node[%s]'s depend node[%s].field[%s] can't be found in candidates[%v]", node.Name, field.Node, field.Name, refers)
			}

			dependNode := g.nodes[field.Node]
			dependNode.AddNext(node.Name) // build edge: prev -> node
			node.AddPrev(field.Node)      // build edge: node -> prev
			node.AddFieldRef(ref)
		}
	}
	return nil
}

// hasLoop returns true if the graph has a loop using topological sort
func (g *Graph) hasLoop() bool {
	// Calculate in-degree of all nodes
	inDegrees := make(map[string]int)
	for name, node := range g.nodes {
		inDegrees[name] = len(node.GetPrevNodes())
	}

	// Find nodes with zero in-degree
	zeroDegreeNodes := make([]string, 0)
	for name, degree := range inDegrees {
		if degree == 0 {
			zeroDegreeNodes = append(zeroDegreeNodes, name)
		}
	}

	// Perform topological sort
	deleteCnt := 0
	for len(zeroDegreeNodes) > 0 {
		node := zeroDegreeNodes[0]
		zeroDegreeNodes = zeroDegreeNodes[1:]
		for _, to := range g.nodes[node].GetNextNodes() {
			inDegrees[to]--
			if inDegrees[to] == 0 {
				zeroDegreeNodes = append(zeroDegreeNodes, to)
			}
		}
		deleteCnt++
	}
	return deleteCnt != len(g.nodes)
}
