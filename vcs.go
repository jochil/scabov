/*
Package vcs
 */
package vcs

import "log"


func Load() {
	repo := &Repository{
		Remote: "https://github.com/jochil/dummy.git",
		System: GIT,
	}

	repo.Init()
	repo.AllCommits()
	repo.AllDevelopers()
	repo.FindCommit("83a5be867e7a2355646aeaa0be1389b54bfa1c94")

	file1 := repo.FindFileInCommit("869f47702d5b1ec4221e4833a008551792a14632", "83a5be867e7a2355646aeaa0be1389b54bfa1c94")
	log.Println(string(file1.Contents[:]))

	file2 := repo.FindFileInCommit("ca564cbf4736a6f5dfeadae9e14303e0c5e1ad3d", "a2f091edd4143080b9b98493485af80de6e9dbff")
	log.Println(string(file2.Contents[:]))

}
