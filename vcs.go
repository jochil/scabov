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
	repo.Commits()
	repo.Developers()
}
