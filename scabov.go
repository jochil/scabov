package main

import (
	"flag"
	"fmt"
	"github.com/jochil/analyzer"
	"github.com/jochil/analyzer/classifier"
	"github.com/jochil/scabov/data"
	"github.com/jochil/vcs"
	"io/ioutil"
	"log"
)

var (
	repoPath = flag.String("p", "", "(remote) path to an vcs repository")
	verbose  = flag.Bool("v", false, "activate verbose output")
	language = flag.String("l", "", "select programming language for analysis")
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

	repo := vcs.NewRepository(*repoPath)

	//TESTING

	//create raw data matrix
	rawMatrix := map[string]map[string]float64{}
	for _, dev := range repo.Developers {
		devData := data.NewDevData(dev)

		rawMatrix[devData.Dev] = map[string]float64{
			"commits":       float64(devData.Commits),
			"files_added":   float64(devData.FileDiff.NumAdded),
			"files_removed": float64(devData.FileDiff.NumRemoved),
			"files_changed": float64(devData.FileDiff.NumChanged),
			//"lines_added":   float64(devData.LineDiff.NumAdded),
			//"lines_removed": float64(devData.LineDiff.NumRemoved),
			//"elements_used": float64(devData.LanguageUsage.NumUsedElements()),
		}

	}

	/*rawMatrix["rama"] = map[string]float64{
		"k": 1,
		"p": 2,
		"v": 1,
	}
	rawMatrix["homa"] = map[string]float64{
		"k": 2,
		"p": 3,
		"v": 3,
	}
	rawMatrix["flora"] = map[string]float64{
		"k": 3,
		"p": 2,
		"v": 1,
	}
	rawMatrix["sb"] = map[string]float64{
		"k": 5,
		"p": 4,
		"v": 7,
	}
	rawMatrix["holl"] = map[string]float64{
		"k": 6,
		"p": 7,
		"v": 6,
	}*/

	/*rawMatrix = map[string]map[string]float64{}
	  rawMatrix["A"] = map[string]float64{
	     "files_added":  33.0 ,
	     "files_removed": 57.0,
	     "files_changed": 10.0,
	  }
	  rawMatrix["B"] = map[string]float64{
	     "files_added":   44.0,
	     "files_removed": 42.0,
	     "files_changed": 21.0,
	  }
	  rawMatrix["C"] = map[string]float64{
	     "files_added":   7.0,
	     "files_removed": 14.0,
	     "files_changed": 2.0,
	  }*/

	//log.Println(rawMatrix)

	matrix := classifier.QCorrelationCoefficient(rawMatrix)
	//matrix := classifier.SquaredEuclideanDistance(rawMatrix)

	/*matrix := map[string]map[string]float64{}
	matrix["kp"] = map[string]float64{
		"kp": 0.0,
		"sp": 8.7,
		"ap": 25.3,
		"lib": 33.7,
		"zp": 37.9,
		"cvp": 49.3,
		"kon": 50.2,
	}

	matrix["sp"] = map[string]float64{
		"kp": 8.7,
		"sp": 0,
		"ap": 14.8,
		"lib": 19.0,
		"zp": 33.2,
		"cvp": 50.5,
		"kon": 40.0,
	}

	matrix["ap"] = map[string]float64{
		"kp": 25.3,
		"sp": 14.8,
		"ap": 0.0,
		"lib": 10.0,
		"zp": 17.8,
		"cvp": 21.3,
		"kon": 24.3,
	}

	matrix["lib"] = map[string]float64{
		"kp": 33.7,
		"sp": 19.0,
		"ap": 10.0,
		"lib": 0,
		"zp": 10.5,
		"cvp": 18.9,
		"kon": 12.9,
	}

	matrix["zp"] = map[string]float64{
		"kp": 37.9,
		"sp": 33.2,
		"ap": 17.8,
		"lib": 10.5,
		"zp": 0,
		"cvp": 7.6,
		"kon": 8.1,
	}

	matrix["cvp"] = map[string]float64{
		"kp": 49.3,
		"sp": 50.5,
		"ap": 21.3,
		"lib": 18.9,
		"zp": 7.6,
		"cvp": 0.0,
		"kon": 7.3,
	}

	matrix["kon"] = map[string]float64{
		"kp": 50.2,
		"sp": 40.0,
		"ap": 24.3,
		"lib": 12.9,
		"zp": 8.1,
		"cvp": 7.3,
		"kon": 0,
	}*/

	log.Println("-----------")
	for x, proximities := range matrix {

		out := x + ": \t\t"
		for y, proximity := range proximities {
			out += fmt.Sprintf("%.3f (%s)\t\t", proximity, y)
		}

		log.Println(out)
	}
	log.Println("-----------")

	groups := classifier.Merge(matrix)

	for i, group := range groups {
		log.Printf("-------Group %d--------", i+1)
		for _, id := range group.Objects {
			log.Println(id, rawMatrix[id])
		}

		log.Println("---------------")
	}

	/*log.Println("-------------")
	  testCommit := repo.Commits["9683bb1198833ffdb7d523702f7adc2b052969f5"]

	  parser := analyzer.NewParser()
	  for _, file := range testCommit.Files {

	     parser.Elements(file)
	  }
	*/

}
