package graph

import (
	"github.com/ongniud/taskflow/model/config"
)

type FieldRef struct {
	Node string
	Idx  int
}

type Node struct {
	*config.Node
	ID   int
	prev map[string]struct{}
	next map[string]struct{}
	refs []*FieldRef
}

func NewNode(nd *config.Node) *Node {
	return &Node{
		Node: nd,
		prev: make(map[string]struct{}),
		next: make(map[string]struct{}),
	}
}

func (n *Node) GetInDegree() int {
	return len(n.prev)
}

func (n *Node) GetPrevNodes() map[string]struct{} {
	return n.prev
}

func (n *Node) GetNextNodes() []string {
	var rs []string
	for s := range n.next {
		rs = append(rs, s)
	}
	return rs
}

func (n *Node) GetFieldRefs() []*FieldRef {
	return n.refs
}

func (n *Node) AddFieldRef(ref *FieldRef) {
	n.refs = append(n.refs, ref)
}

func (n *Node) AddPrev(nd string) {
	n.prev[nd] = struct{}{}
}

func (n *Node) AddNext(nd string) {
	n.next[nd] = struct{}{}
}
func (n *Node) GetParam() string {
	return n.Param
}
