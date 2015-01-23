package vcs

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	_ = iota
	GIT
	SVN
)

//internal representation of an repository
type Repository struct {
	Commits    map[string]*Commit
	Developers map[string]*Developer

	path      string
	Workspace string
}

/*
initialize the repository (connector)
*/
func NewRepository(path string) (*Repository, error) {

	//TODO evaluate vcs type
	system := GIT

	repo := &Repository{
		path: path,
	}

	repo.checkWorkspace()

	//get correct connector for given vcs
	var connector Connector
	switch system {
	case GIT:
		connector = &GitConnector{}
	}

	//local or remote path?
	if _, err := os.Stat(path); err == nil {
		if err := connector.LoadLocal(path, repo.Workspace); err != nil {
			return nil, err
		}
	} else if err := connector.LoadRemote(path, repo.Workspace); err != nil {
		return nil, err
	}

	repo.Commits = connector.Commits()
	repo.Developers = connector.Developers()

	return repo, nil
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
			log.Printf("found file %s in commit %s", file.Id, commitId)
			return file
		}
	}
	return nil
}

//TODO validate directories
func (r *Repository) checkWorkspace() {
	if r.Workspace == "" {
		//get hash from repo url
		h := sha1.New()
		io.WriteString(h, r.path)
		dir := fmt.Sprintf("%x", h.Sum(nil))

		//get current directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(cwd)
		}
		//TODO ensure file/dir handling
		r.Workspace = filepath.Join(cwd, "workspace", dir)
	}
}
