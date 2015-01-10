package main

import (
	"flag"
	"github.com/jochil/analyzer"
	"github.com/jochil/analyzer/classifier"
	"github.com/jochil/scabov/export"
	"github.com/jochil/vcs"
	"io/ioutil"
	"log"
	"os"
)

var (
	repoPath       = flag.String("p", "", "(remote) path to an vcs repository")
	verbose        = flag.Bool("v", false, "activate verbose output")
	language       = flag.String("l", "", "select programming language for analysis")
	metrics        = flag.Bool("m", false, "activate metrics calculation")
	classification = flag.Bool("c", false, "activate developer classification")
	outputFilename = flag.String("o", "result.xml", "select output file")

	outputFile *os.File
)

func main() {

	flag.Parse()

	// setup logging
	log.SetFlags(0)
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	filter := vcs.NewLanguageFilter(*language)
	vcs.Filter = filter
	analyzer.Filter = filter

	//TODO path validation?
	outputFile, _ = os.Create(*outputFilename)

	// load repo
	if *repoPath == "" {
		log.Fatal("repository path missing, e.g.: -p \"mypath/repo\"")
	}

	repo := vcs.NewRepository(*repoPath)

	if *classification {
		executeClassification(repo)
	}

	log.Printf("Saved results to %s", *outputFilename)
	export.SaveFile(outputFile)

	log.Println("Finsihed")
}

func executeClassification(repo *vcs.Repository) {

	log.Println("start classification")

	//create raw data matrix
	rawMatrix := map[string]map[string]float64{}
	for _, dev := range repo.Developers {

		cycloDiff := analyzer.CalcCycloDiff(dev)
		fileDiff := dev.FileDiff()
		lineDiff := dev.LineDiff()

		//remove dev without any contribution
		if cycloDiff.IsEmpty() && fileDiff.IsEmpty() && lineDiff.IsEmpty() {
			continue;
		}

		rawMatrix[dev.Id] = map[string]float64{
			"files_added":     float64(fileDiff.Added),
			"files_removed":   float64(fileDiff.Removed),
			"files_changed":   float64(fileDiff.Changed),
			"lines_added":     float64(lineDiff.Added),
			"lines_removed":   float64(lineDiff.Removed),
			"cyclo_increased": float64(cycloDiff.Increased),
			"cyclo_decreased": float64(cycloDiff.Decreased),
		}

	}
	log.Println("\t created raw data matrix")
	//export.PrintMatrix(rawMatrix)

	matrix := classifier.QCorrelationCoefficient(rawMatrix)
	log.Println("\t calculated distance matrix")
	//export.PrintMatrix(matrix)

	groups := classifier.Merge(matrix)
	log.Printf("\t finished classification, found %d groups within %d active devs", len(groups), len(rawMatrix))

	export.SaveClassificationResult(groups, rawMatrix)
}
