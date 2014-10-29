/*
Package vcs
 */
package vcs

func LoadRepository(path string) *Repository {

	repo := &Repository{
		Remote: path,
		System: GIT,
	}
	repo.Init()

	return repo
}
