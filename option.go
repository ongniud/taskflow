package taskflow

import (
	"github.com/ongniud/taskflow/model/graph"
)

const (
	// DefaultPoolSize is the default pool size.
	DefaultPoolSize = 10
)

type TaskOptions struct {
	Graph       *graph.Graph
	Parallelism int
	Timeout     int
	Debug       bool
}

type TaskOption func(*TaskOptions)

func WithGraph(g *graph.Graph) TaskOption {
	return func(o *TaskOptions) {
		o.Graph = g
	}
}

func WithTimeout(timeout int) TaskOption {
	return func(o *TaskOptions) {
		o.Timeout = timeout
	}
}

func WithParallelism(parallel int) TaskOption {
	return func(o *TaskOptions) {
		o.Parallelism = parallel
	}
}

func WithDebug(debug bool) TaskOption {
	return func(o *TaskOptions) {
		o.Debug = debug
	}
}
