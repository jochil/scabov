package analyzer

import (
	"github.com/gyuho/goraph/graph/gs"
	"log"
	"os/exec"
)

/*
calculates McCabe-number/cyclomatic compexitiy of the current cfg
formula: M = e - n + 2p
e: number of edges
n: number of nodes
p: number of nodes without outgoing edge
*/
func CyclomaticComplexity(cfg *gs.Graph) int {

	e := cfg.GetEdgesSize()
	n := cfg.GetVerticesSize()
	p := 1

	return e - n + 2*p
}

func dumpCFG(name string, cfg *gs.Graph) {
	path := "workspace/graphs/" + name
	cfg.ToDOTFile(path + ".dot")
	cmd := exec.Command("dot", "-Tpng", path+".dot", "-o", path+".png")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
