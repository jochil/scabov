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
	"path"
)

var (
	//parameters
	repoPath       = flag.String("p", "", "(remote) path to an vcs repository")
	verbose        = flag.Bool("v", false, "activate verbose output")
	language       = flag.String("l", "", "select programming language for analysis")
	metrics        = flag.Bool("m", false, "activate metrics calculation")
	classification = flag.Bool("c", false, "activate developer classification")
	outputFilename = flag.String("o", "", "select output file")

	//local vars
	repo                                                  *vcs.Repository
	outputFile                                            *os.File
	styleGroups, contributionGroups                       []*classifier.Group
	runStyleClassification, runContributionClassification bool
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

	// load repo
	if *repoPath == "" {
		log.Fatal("repository path missing, e.g.: -p \"mypath/repo\"")
	}

	var err error
	repo, err = vcs.NewRepository(*repoPath)
	if err != nil {
		log.Fatal(err)
	}

	//TODO path validation?
	if *outputFilename == "" {
		outputFile, _ = os.Create(path.Join(repo.Workspace, "result.xml"))
	} else {
		outputFile, _ = os.Create(*outputFilename)
	}

	runStyleClassification = true
	runContributionClassification = true

	if *classification {
		executeCompleteClassification()
	}

	if *metrics {
		executeMetricsCalculation()
	}

	log.Printf("saved results to %s", outputFile.Name())
	export.SaveFile(outputFile)

	//TODO clean up (delete workspace, ...)

	log.Println("finsihed")
}

func executeMetricsCalculation() {

	if runStyleClassification == true {
		executeStyleClassification()
	}
	if runContributionClassification == true {
		executeContributionClassification()
	}

	log.Println("started metric extraction")

	styleHomogeneity := analyzer.CalcHomogeneity(styleGroups)
	log.Printf("\t style homogeneity: %.2f", styleHomogeneity)

	contributionHomogeneity := analyzer.CalcHomogeneity(contributionGroups)
	log.Printf("\t contribution homogeneity: %.2f", contributionHomogeneity)

	analyzer.LoadHistory(repo)

	stability := analyzer.CalcFunctionStability(repo)
	log.Printf("\t overall function stability: %.2f", stability)

	export.SaveMetricsResult(stability, styleHomogeneity, contributionHomogeneity)
	export.SaveFunctions(analyzer.History)
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
	runStyleClassification = false
}

func executeContributionClassification() {
	log.Println("started contribution classification")
	contributionRawMatrix := analyzer.ContributionData(repo)
	contributionGroups = classifier.ClusterAnalysis(contributionRawMatrix)
	export.SaveClassificationResult("contribution", contributionGroups, contributionRawMatrix)
	log.Println("\t finished contribution classification")
	runContributionClassification = false
}
