package main

import (
	"github.com/jochil/vcs"
	"log"
	"flag"
	"io/ioutil"
)

var (
	repoPath = flag.String("p", "", "(remote) path to an vcs repository")
	verbose  = flag.Bool("v", false, "activate verbose output")
)

func main() {

	flag.Parse()

	// setup logging
	log.SetFlags(0)
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	// load repo
	if *repoPath == "" {
		log.Fatal("repository path missing, e.g.: -p \"mypath/repo\"")
	}

	repo := vcs.LoadRepository(*repoPath)

	firstCommit := repo.FirstCommit()
	log.Println("First commit", firstCommit)

	/*
	repo.AllDevelopers()
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
