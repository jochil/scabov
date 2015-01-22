package analyzer

import (
	vcs "github.com/jochil/scabov/vcs"
)

var Filter vcs.LanguageFilter

func ContributionData(repo *vcs.Repository) (rawData map[string]map[string]float64) {

	rawData = map[string]map[string]float64{}

	for _, dev := range repo.Developers {

		complexityDiff := CalcComplexityDiff(dev)
		fileDiff := dev.FileDiff()
		lineDiff := dev.LineDiff()

		if fileDiff.IsEmpty() == false ||
			lineDiff.IsEmpty() == false ||
			complexityDiff.CycloIncreased > 0 ||
			complexityDiff.CycloDecreased > 0 {

			rawData[dev.Id] = map[string]float64{
				"files_added":     float64(fileDiff.Added),
				"files_removed":   float64(fileDiff.Removed),
				"files_changed":   float64(fileDiff.Changed),
				"lines_added":     float64(lineDiff.Added),
				"lines_removed":   float64(lineDiff.Removed),
				"cyclo_increased": float64(complexityDiff.CycloIncreased),
				"cyclo_decreased": float64(complexityDiff.CycloDecreased),
			}
		}
	}

	//export.PrintMatrix(rawData)
	return rawData
}

func StyleData(repo *vcs.Repository) (rawData map[string]map[string]float64) {

	rawData = map[string]map[string]float64{}

	overallCycloMax := 0
	overallFuncNodesMax := 0

	for _, dev := range repo.Developers {

		complexityDiff := CalcComplexityDiff(dev)
		languageUsage := CalcLanguageUsage(dev)

		if complexityDiff.CycloAvg() != 0.0 ||
			languageUsage.Value() != 0.0 ||
			complexityDiff.FuncNodesAvg() != 0.0 {

			if funcNodesMax := complexityDiff.FuncNodesMax(); funcNodesMax > overallFuncNodesMax {
				overallFuncNodesMax = funcNodesMax
			}
			if cycloMax := complexityDiff.CycloMax(); cycloMax > overallCycloMax {
				overallCycloMax = cycloMax
			}

			rawData[dev.Id] = map[string]float64{
				"cyclo_avg":      complexityDiff.CycloAvg(),
				"language_usage": languageUsage.Value(),
				"function_size":  complexityDiff.FuncNodesAvg(),
			}
		}
	}

	//normalize data
	for _, row := range rawData {
		crtCyclo := row["cyclo_avg"]
		row["cyclo_avg"] = crtCyclo * 100.0 / float64(overallCycloMax)

		crtFuncSize := row["function_size"]
		row["function_size"] = crtFuncSize * 100.0 / float64(overallFuncNodesMax)
	}

	//export.PrintMatrix(rawData)
	return rawData
}
