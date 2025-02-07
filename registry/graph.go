package registry

import (
	"errors"
	"sync"

	"github.com/ongniud/taskflow/model/graph"
)

var (
	graphs sync.Map
)

func RegisterGraph(g *graph.Graph) error {
	if g == nil || g.Name == "" {
		return errors.New("empty graph")
	}
	graphs.Store(g.Name, g)
	return nil
}

func GetGraph(gr string) *graph.Graph {
	gi, ok := graphs.Load(gr)
	if !ok {
		return nil
	}
	g, ok := gi.(*graph.Graph)
	if !ok {
		return nil
	}
	return g
}
