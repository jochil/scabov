package vcs

import (
	"fmt"
	"io"
	"crypto/sha1"
	"path/filepath"
	"os"
	"log"
)

const (
	_ = iota
	GIT
	SVN
)

//internal representation of an repository
type Repository struct {
	Commits      map[string]*Commit
	Developers   map[string]*Developer

	connector    Connector
	remote       string
	local        string
}

/*
initialize the repository (connector)
 */
func NewRepository(path string) *Repository {

	//TODO evaluate vcs type
	system := GIT

	repo := &Repository{
		remote: path,
	}

	repo.checkWorkspace()

	//get correct connector for given system
	switch system {
	case GIT:
		repo.connector = &GitConnector{}
	}

	repo.Commits, repo.Developers = repo.connector.Load(repo.remote, repo.local)

	return repo
}

//TODO replace this naive approach
func (r *Repository) FirstCommit() *Commit {
	for _, commit := range r.Commits {
		if len(commit.Parents) == 0 {
			return commit
		}
	}
	return nil
}

//Lookup for a single commit
func (r *Repository) FindCommit(id string) *Commit {
	//TODO implement lockup without questioning all commits
	commits := r.Commits
	if commit, ok := commits[id]; ok {
		log.Printf("found commit with id %s", id)
		return commit
	} else {
		return nil
	}
}

//Looks for a specific file in a given commit
func (r *Repository) FindFileInCommit(fileId string, commitId string) *File {
	commit := r.FindCommit(commitId)
	if commit != nil {
		if file, ok := commit.Files[fileId]; ok {
			log.Printf("found file %s (%s) in commit %s", file.Id, file.Path, commitId)
			return file
		}
	}
	return nil
}

//TODO validate directories
func (r *Repository) checkWorkspace() {
	if r.local == "" {
		//get hash from repo url
		h := sha1.New()
		io.WriteString(h, r.remote)
		dir := fmt.Sprintf("%x", h.Sum(nil))

		//get current directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(cwd)
		}
		//TODO ensure file/dir handling
		r.local = filepath.Join(cwd, "workspace", dir)
	}
}
