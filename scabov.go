package main

import (
	"github.com/jochil/vcs"
	"log"
	"flag"
	"io/ioutil"
	"github.com/jochil/analyzer"
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



	log.Println("-------------")
	for _, dev := range repo.Developers {
		for _, commit := range dev.Commits {
			log.Println(commit)
			for path, file := range commit.Files {
				log.Println("\t", path, file)
			}
		}
	}

	/*
	log.Println("-------------")
	for _, dev := range repo.Developers {
		devData := data.DevData{}
		devData.Dev = dev.Id
		devData.Commits = uint16(len(dev.Commits))
		devData.LineDiff = analyzer.CalcLineDiff(dev)
		log.Println(devData.String())
		log.Printf("\t FileDiff: Added: %d, Removed: %d, Changed: %d",
			len(dev.NewFiles), len(dev.RemovedFiles), len(dev.ChangedFiles),
		)
	}

	*/
	//firstCommit := repo.FirstCommit()
	//log.Println("First commit", firstCommit)

	/*

	file1 := repo.FindFileInCommit("869f47702d5b1ec4221e4833a008551792a14632", "83a5be867e7a2355646aeaa0be1389b54bfa1c94")
	file2 := repo.FindFileInCommit("ca564cbf4736a6f5dfeadae9e14303e0c5e1ad3d", "a2f091edd4143080b9b98493485af80de6e9dbff")

	log.Println(analyzer.CountLineDiff(file1.String(), file2.String()))

	repo.AllCommits()

	repo.FindCommit("83a5be867e7a2355646aeaa0be1389b54bfa1c94")

	file1 := repo.FindFileInCommit("869f47702d5b1ec4221e4833a008551792a14632", "83a5be867e7a2355646aeaa0be1389b54bfa1c94")
	log.Println(string(file1.Contents[:]))

	file2 := repo.FindFileInCommit("ca564cbf4736a6f5dfeadae9e14303e0c5e1ad3d", "a2f091edd4143080b9b98493485af80de6e9dbff")
	log.Println(string(file2.Contents[:]))
	*/
}
