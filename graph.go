package analyzer

import (
	"fmt"
	"log"
)

// counter for node id's
var lastNodeID uint = 0

// node struct
type Node struct {
	ID    uint
	Label string
}

// constructor
func NewNode(label string) *Node {

	lastNodeID++
	node := &Node{
		ID: lastNodeID,
		Label: label,
	}
	return node
}

func (node *Node) String() string {
	return fmt.Sprintf("Node %d (%s)", node.ID, node.Label)
}

// edge struct
type Edge struct {
	Start *Node
	End   *Node
}

// control flow graph
type ControlFlowGraph struct {
	nodes map[uint]*Node
	edges []*Edge
}

// constructor
func NewControlFlowGraph() *ControlFlowGraph {
	cfg := &ControlFlowGraph{
		nodes: map[uint]*Node{},
		edges: make([]*Edge, 0),
	}
	return cfg
}

/*
calculates McCabe-number/cyclomatic compexitiy of the current cfg
formula: M = e - n + 2p
e: number of edges
n: number of nodes
p: number of nodes without outgoing edge
*/
func (cfg *ControlFlowGraph) CyclomaticComplexity() int {

	//calculate nodes without outgoing edge
	endNodes := map[uint]*Node{}
	//-- clone map with all nodes...
	for k, v := range cfg.nodes {
		endNodes[k] = v
	}
	//... and remove all who are obviously no end-node
	for _, edge := range cfg.edges {
		delete(endNodes, edge.Start.ID)
	}

	e := len(cfg.edges)
	n := len(cfg.nodes)
	p := len(endNodes)


	return e - n + 2 * p
}

// adds a node to the control flow graph
func (cfg *ControlFlowGraph) Add(node *Node) {
	cfg.nodes[node.ID] = node
}

// creates an edge between to nodes
func (cfg *ControlFlowGraph) Connect(startNode *Node, endNode *Node) {
	edge := &Edge{Start: startNode, End: endNode}
	cfg.edges = append(cfg.edges, edge)
}

func (cfg *ControlFlowGraph) String() string {
	return fmt.Sprintf("CFG Nodes: %d, Edges: %d", len(cfg.nodes), len(cfg.edges))
}

// prints the whole control flow graph
func (cfg *ControlFlowGraph) Print() {

	log.Println(cfg.String())
	log.Printf("cyclo: %d", cfg.CyclomaticComplexity())

	for _, node := range cfg.nodes {
		log.Printf("%s", node.String())
	}

	log.Println("-----------------")

	for _, edge := range cfg.edges {
		log.Printf("%s -> %s", edge.Start.String(), edge.End.String())
	}
}
