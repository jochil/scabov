package main

import (
	"flag"
	"github.com/jochil/analyzer"
	"github.com/jochil/analyzer/classifier"
	"github.com/jochil/scabov/data"
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

	log.Println("finsihed")
}

func executeClassification(repo *vcs.Repository) {

	log.Println("start classification")

	//create raw data matrix
	rawMatrix := map[string]map[string]float64{}
	for _, dev := range repo.Developers {
		devData := data.NewDevData(dev)

		rawMatrix[devData.Dev] = map[string]float64{
			"commits": float64(devData.Commits),
			//"files_added":   float64(devData.FileDiff.NumAdded),
			//"files_removed": float64(devData.FileDiff.NumRemoved),
			//"files_changed": float64(devData.FileDiff.NumChanged),
			"lines_added":   float64(devData.LineDiff.NumAdded),
			"lines_removed": float64(devData.LineDiff.NumRemoved),
			//"elements_used": float64(devData.LanguageUsage.NumUsedElements()),
		}

	}
	log.Println("\tcreated raw data matrix")

	matrix := classifier.QCorrelationCoefficient(rawMatrix)
	log.Println("\tscalculated distance matrix")

	groups := classifier.Merge(matrix)
	log.Printf("\tfinished classification, found %d groups within %d devs", len(groups), len(rawMatrix))

	export.SaveClassificationResult(groups, rawMatrix)
}
