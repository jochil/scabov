package analyzer

import "github.com/jochil/scabov/analyzer/classifier"

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
