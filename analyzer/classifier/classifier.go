package classifier

import "log"

func ClusterAnalysis(rawMatrix map[string]map[string]float64) []*Group {

	distanceMatrix := QCorrelationCoefficient(rawMatrix)
	log.Println("\t calculated distance matrix")
	//export.PrintMatrix(distanceMatrix)

	groups := Merge(distanceMatrix)
	log.Printf("\t finished classification, found %d groups within %d pattern",
		len(groups), len(distanceMatrix))

	return groups
}
