package main

import (
	"flag"
	"github.com/jochil/analyzer"
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
	log.Println("-------------")
	for _, dev := range repo.Developers {
		log.Println(data.NewDevData(dev))
	}
}
