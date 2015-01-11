package analyzer

import (
	"github.com/gyuho/goraph/graph/gs"
	"github.com/jochil/scabov/vcs"
	"log"
	"os/exec"
	"math"
)

type CycloDiff struct {
	Increased int
	Decreased int
	New       []int
}

func (diff *CycloDiff) IsEmpty() bool {
	return diff.Increased == 0 && diff.Decreased == 0 && len(diff.New) == 0
}

func (diff *CycloDiff) Sum() int {
	sum := 0
	for _, value := range diff.New {
		sum += value
	}
	return sum
}

func (diff *CycloDiff) Avg() float64 {
	if result := float64(diff.Sum()) / float64(len(diff.New)); math.IsNaN(result) == false {
		return result
	}
	return 0
}

func (diff *CycloDiff) Max() int {
	max := 0
	for _, value := range diff.New {
		if value > max {
			max = value
		}
	}
	return max
}

func CalcCycloDiff(dev *vcs.Developer) CycloDiff {

	parser := NewParser()

	cycloDiff := CycloDiff{0, 0, []int{}}

	//handle added files
	for _, file := range dev.AddedFiles() {
		functions := parser.Functions(file)

		for _, function := range functions {
			cyclo := CyclomaticComplexity(function.CFG)
			cycloDiff.New = append(cycloDiff.New, cyclo)
		}
	}

	//handle modified files
	for _, file := range dev.ModifiedFiles() {

		//TODO just handle one parent file, get this working for n-parents
		if parentFile := file.Parents[0]; parentFile != nil {

			functions := parser.Functions(file)
			parentFunctions := parser.Functions(parentFile)

			for name, function := range parentFunctions {

				newCyclo := CyclomaticComplexity(function.CFG)

				if parentFunction, ok := functions[name]; ok {
					oldCyclo := CyclomaticComplexity(parentFunction.CFG)

					if oldCyclo > newCyclo {
						cycloDiff.Decreased++
					} else if oldCyclo < newCyclo {
						cycloDiff.Increased++
					}

				} else {
					cycloDiff.New = append(cycloDiff.New, newCyclo)
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
