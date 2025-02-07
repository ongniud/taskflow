package graph

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ongniud/taskflow/model/config"
)

func TestGraphParsing(t *testing.T) {
	g := &config.Graph{
		Name: "graph",
		Nodes: []*config.Node{
			{
				Name:     "start",
				Kind:     config.NodeKindOperator,
				Operator: "debug",
				Inputs:   []*config.Field{},
				Outputs: []*config.Field{
					{Name: "s1"},
					{Name: "s2"},
					{Name: "s3"},
				},
			},
			{
				Name:     "middle-a",
				Kind:     config.NodeKindOperator,
				Operator: "debug",
				Inputs: []*config.Field{
					{Name: "s1"},
					{Name: "s2"},
					{Name: "s3"},
				},
				Outputs: []*config.Field{
					{Name: "m1"},
					{Name: "m2"},
					{Name: "m3"},
				},
			},
			{
				Name:     "middle-b",
				Kind:     config.NodeKindOperator,
				Operator: "debug",
				Inputs: []*config.Field{
					{Name: "s1"},
					{Name: "s2"},
					{Name: "s3"},
				},
				Outputs: []*config.Field{
					{Name: "m4"},
					{Name: "m5"},
					{Name: "m6"},
				},
			},
			{
				Name:     "end",
				Kind:     config.NodeKindOperator,
				Operator: "debug",
				Inputs: []*config.Field{
					{Name: "m1"},
					{Name: "m2"},
					{Name: "m3"},
					{Name: "m4"},
					{Name: "m5"},
					{Name: "m6"},
				},
				Outputs: []*config.Field{
					{Name: "e1"},
					{Name: "e2"},
					{Name: "e3"},
				},
			},
		},
	}

	_, err := NewGraph(g)
	if err != nil {
		t.Errorf("Failed to parse graph: %v", err)
		return
	}

	graphStr, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("graph:", string(graphStr))
}
