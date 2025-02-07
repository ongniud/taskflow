package config

type Field struct {
	Node    string `json:"node,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Require bool   `json:"require,omitempty"`
	Mapping string `json:"mapping,omitempty"`
}

const (
	NodeKindOperator = "operator"
	NodeKindGraph    = "graph"
)

type Node struct {
	// Base
	Name      string `json:"name,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Operator  string `json:"operator,omitempty"`
	Graph     string `json:"graph,omitempty"`
	Condition string `json:"condition,omitempty"`

	Inputs  []*Field `json:"inputs,omitempty"`
	Outputs []*Field `json:"outputs,omitempty"`
	Param   string   `json:"param,omitempty"`

	//params  []*Field               `json:"params" yaml:"params"`
	//Args    map[string]interface{} `json:"args" yaml:"args"`

	// Control
	Timeout   int64 `json:"timeout,omitempty"`
	ErrIgnore bool  `json:"err_ignore,omitempty"`
	ErrPrune  bool  `json:"err_prune,omitempty"`
	ErrAbort  bool  `json:"err_abort,omitempty"`
}
