package main

import (
	"flag"
	"github.com/jochil/scabov/analyzer"
	"github.com/jochil/scabov/analyzer/classifier"
	"github.com/jochil/scabov/export"
	"github.com/jochil/scabov/vcs"
	"io/ioutil"
	"log"
	"os"
)

var (
	//parameters
	repoPath       = flag.String("p", "", "(remote) path to an vcs repository")
	verbose        = flag.Bool("v", false, "activate verbose output")
	language       = flag.String("l", "", "select programming language for analysis")
	metrics        = flag.Bool("m", false, "activate metrics calculation")
	classification = flag.Bool("c", false, "activate developer classification")
	outputFilename = flag.String("o", "result.xml", "select output file")

	//local vars
	repo                            *vcs.Repository
	outputFile                      *os.File
	styleGroups, contributionGroups []*classifier.Group
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

	var err error
	repo, err = vcs.NewRepository(*repoPath)
	if err != nil {
		log.Fatal(err)
	}


	if *classification {
		executeCompleteClassification()
	}

	if *metrics {
		executeMetricsCalculation()
	}

	log.Printf("Saved results to %s", *outputFilename)
	export.SaveFile(outputFile)

	log.Println("Finsihed")
}

func executeMetricsCalculation() {

	if styleGroups == nil {
		executeStyleClassification()
	}
	if contributionGroups == nil {
		executeContributionClassification()
	}

	log.Println("started metric extraction")
}

func executeCompleteClassification() {
	executeStyleClassification()
	executeContributionClassification()
}

func executeStyleClassification() {
	log.Println("started style classification")
	styleRawMatrix := analyzer.StyleData(repo)
	styleGroups = classifier.ClusterAnalysis(styleRawMatrix)
	export.SaveClassificationResult("style", styleGroups, styleRawMatrix)
	log.Println("\t finished style classification")
}

func executeContributionClassification() {
	log.Println("started contribution classification")
	contributionRawMatrix := analyzer.ContributionData(repo)
	contributionGroups = classifier.ClusterAnalysis(contributionRawMatrix)
	export.SaveClassificationResult("contribution", contributionGroups, contributionRawMatrix)
	log.Println("\t finished contribution classification")
}
