package analyzer

import (
	"github.com/jochil/scabov/analyzer/classifier"
	"github.com/jochil/scabov/vcs"
)

func CalcHomogeneity(groups []*classifier.Group) float64 {

	if len(groups) == 1 {
		return 1.0
	}

	if len(groups) == 0 {
		return 0.0
	}

	countDevs := 0
	for _, group := range groups {
		countDevs += len(group.Objects)
	}

	return float64(countDevs-len(groups)) / float64(countDevs)
}

func CalcFunctionStability(repo *vcs.Repository) float64 {

	count := 0.0
	sum := 0.0
	for _, fileHistory := range History {
		for _, history := range fileHistory {
			count++
			sum += history.Stability()
		}
	}

	stability := sum / count

	return stability
}
