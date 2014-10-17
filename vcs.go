/*
Package vcs
 */
package vcs

func Load() {
	repo := &Repository{
		remote: "https://github.com/humhub/humhub.git",
		system: GIT,
	}

	repo.Init()
	repo.Commits()
	repo.Developers()
}
