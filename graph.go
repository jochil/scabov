package analyzer

import (
	"github.com/gyuho/goraph/graph/gs"
	"github.com/jochil/vcs"
	"log"
	"os/exec"
)

type CycloDiff struct {
	Increased int
	Decreased int
}

func (diff *CycloDiff) IsEmpty() bool {
	return diff.Increased == 0 && diff.Decreased == 0
}

func CalcCycloDiff(dev *vcs.Developer) CycloDiff {

	parser := NewParser()

	cycloDiff := CycloDiff{0, 0}

	for _, file := range dev.ModifiedFiles() {

		//TODO just handle one parent file, get this working for n-parents
		if parentFile := file.Parents[0]; parentFile != nil {

			functions := parser.Functions(file)
			parentFunctions := parser.Functions(parentFile)

			for parentName, parentFunction := range parentFunctions {

				if function, ok := functions[parentName]; ok {
					newCyclo := CyclomaticComplexity(function.CFG)
					oldCyclo := CyclomaticComplexity(parentFunction.CFG)

					if oldCyclo > newCyclo {
						cycloDiff.Decreased++
					} else if oldCyclo < newCyclo {
						cycloDiff.Increased++
					}
				}
			}
		}
	}

	return cycloDiff
}

/*
calculates McCabe-number/cyclomatic compexitiy of the current cfg
formula: M = e - n + p
e: number of edges
n: number of nodes
p: number of connected components
*/
func CyclomaticComplexity(cfg *gs.Graph) int {
	e := cfg.GetEdgesSize()
	n := cfg.GetVerticesSize()
	p := 1
	return e - n + p
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
