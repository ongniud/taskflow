package config

type Graph struct {
	Name        string  `json:"name,omitempty"`
	Nodes       []*Node `json:"nodes,omitempty"`
	Timeout     int32   `json:"timeout,omitempty"`
	Parallelism int     `json:"parallelism,omitempty"`
}
