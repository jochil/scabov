/*
Package vcs
 */
package vcs


func Load() {
	repo := &Repository{
		remote: "https://github.com/jochil/dummy.git",
		system: GIT,
	}

	repo.Init()
	//repo.AllCommits()
	//repo.AllDevelopers()
	//repo.FindCommit("83a5be867e7a2355646aeaa0be1389b54bfa1c94")
	repo.FindFileInCommit("869f47702d5b1ec4221e4833a008551792a14632", "83a5be867e7a2355646aeaa0be1389b54bfa1c94")


}
