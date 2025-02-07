package utils

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ongniud/taskflow/model/config"
	"github.com/ongniud/taskflow/model/graph"
)

func LoadGraph(path string) (*graph.Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := &config.Graph{}
	if err = json.Unmarshal(content, cfg); err != nil {
		return nil, err
	}

	gra, err := graph.NewGraph(cfg)
	if err != nil {
		return nil, err
	}

	return gra, nil
}
