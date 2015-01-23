package analyzer

import (
	"github.com/gyuho/goraph/graph/gs"
	"github.com/jochil/scabov/vcs"
	"math"
)

type ComplexityDiff struct {
	CycloIncreased int
	CycloDecreased int
	CycloNew       []int
	FuncNodes      []int
}

func (diff *ComplexityDiff) CycloSum() int {
	sum := 0
	for _, value := range diff.CycloNew {
		sum += value
	}
	return sum
}

func (diff *ComplexityDiff) CycloAvg() float64 {
	if result := float64(diff.CycloSum()) / float64(len(diff.CycloNew)); math.IsNaN(result) == false {
		return result
	}
	return 0
}

func (diff *ComplexityDiff) CycloMax() int {
	max := 0
	for _, value := range diff.CycloNew {
		if value > max {
			max = value
		}
	}
	return max
}

func (diff *ComplexityDiff) FuncNodesSum() int {
	sum := 0
	for _, value := range diff.FuncNodes {
		sum += value
	}
	return sum
}

func (diff *ComplexityDiff) FuncNodesAvg() float64 {
	if result := float64(diff.FuncNodesSum()) / float64(len(diff.FuncNodes)); math.IsNaN(result) == false {
		return result
	}
	return 0
}

func (diff *ComplexityDiff) FuncNodesMax() int {
	max := 0
	for _, value := range diff.FuncNodes {
		if value > max {
			max = value
		}
	}
	return max
}

func CalcComplexityDiff(dev *vcs.Developer) ComplexityDiff {

	parser := NewParser()

	diff := ComplexityDiff{0, 0, []int{}, []int{}}

	//handle added files
	for _, file := range dev.AddedFiles() {
		functions := parser.Functions(file)

		for _, function := range functions {
			cyclo := CyclomaticComplexity(function.CFG)
			diff.CycloNew = append(diff.CycloNew, cyclo)
			diff.FuncNodes = append(diff.FuncNodes, function.NumNodes)
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
						diff.CycloDecreased++
					} else if oldCyclo < newCyclo {
						diff.CycloIncreased++
					}

				} else {
					diff.CycloNew = append(diff.CycloNew, newCyclo)
					diff.FuncNodes = append(diff.FuncNodes, function.NumNodes)
				}
			}
		}
	}

	return diff
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
