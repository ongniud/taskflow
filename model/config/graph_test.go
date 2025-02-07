package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseGraph(t *testing.T) {
	jsonData := `{
	  "name": "example_graph",
	  "dag": [
	    {
	      "name": "node1",
	      "op": "operation1",
	      "args": {
	        "arg1": "value1",
	        "arg2": "value2"
	      },
	      "params": [
	        {
	          "node": "node1",
	          "name": "param1",
	          "type": "string",
	          "require": true,
	          "mapping": ""
	        },
	        {
	          "node": "node1",
	          "name": "param2",
	          "type": "int",
	          "require": false,
	          "mapping": ""
	        }
	      ],
	      "inputs": [
	        {
	          "node": "node2",
	          "name": "input1",
	          "type": "string",
	          "require": true,
	          "mapping": ""
	        },
	        {
	          "node": "node3",
	          "name": "input2",
	          "type": "float",
	          "require": false,
	          "mapping": ""
	        }
	      ],
	      "outputs": [
	        {
	          "node": "node4",
	          "name": "output1",
	          "type": "bool",
	          "require": true,
	          "mapping": ""
	        }
	      ],
	      "err_ignore": false,
	      "err_prune": true,
	      "err_abort": false,
	      "timeout": 300
	    },
	    {
	      "name": "node2",
	      "op": "operation2",
	      "args": {},
	      "params": [],
	      "inputs": [],
	      "outputs": [],
	      "err_ignore": false,
	      "err_prune": false,
	      "err_abort": false,
	      "timeout": 0
	    },
	    {
	      "name": "node3",
	      "op": "operation3",
	      "args": {},
	      "params": [],
	      "inputs": [],
	      "outputs": [],
	      "err_ignore": true,
	      "err_prune": false,
	      "err_abort": true,
	      "timeout": 0
	    },
	    {
	      "name": "node4",
	      "op": "operation4",
	      "args": {},
	      "params": [],
	      "inputs": [],
	      "outputs": [],
	      "err_ignore": false,
	      "err_prune": false,
	      "err_abort": false,
	      "timeout": 0
	    }
	  ],
	  "timeout": 600
	}`

	var graph Graph
	err := json.Unmarshal([]byte(jsonData), &graph)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	fmt.Printf("g Name: %s\n", graph.Name)
	fmt.Printf("g Timeout: %d\n", graph.Timeout)

	for _, node := range graph.Nodes {
		fmt.Printf("\nNode Name: %s\n", node.Name)
		fmt.Printf("Node Operation: %s\n", node.Operator)
		fmt.Printf("Node Param: %+v\n", node.Param)

		fmt.Println("Node inputs:")
		for _, input := range node.Inputs {
			fmt.Printf("- Name: %s, Type: %s, Require: %t\n", input.Name, input.Type, input.Require)
		}

		fmt.Println("Node outputs:")
		for _, output := range node.Outputs {
			fmt.Printf("- Name: %s, Type: %s, Require: %t\n", output.Name, output.Type, output.Require)
		}

		fmt.Printf("Node Error Ignore: %t\n", node.ErrIgnore)
		fmt.Printf("Node Error Prune: %t\n", node.ErrPrune)
		fmt.Printf("Node Error Abort: %t\n", node.ErrAbort)
		fmt.Printf("Node Timeout: %d\n", node.Timeout)
	}
}
