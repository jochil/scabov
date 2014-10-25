/*
Package vcs
 */
package vcs

func LoadRemoteRepository(path string) *Repository {

	repo := &Repository{
		Remote: path,
		System: GIT,
	}
	repo.Init()

	return repo
}
