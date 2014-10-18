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
	GIT = 1
)

//internal representation of an repository
type Repository struct {
	remote       string
	local        string
	system       int
	connector    Connector
}

/*
initialize the repository (connector)
TODO should be replaced by a constructor
 */
func (r *Repository) Init() {

	r.checkWorkspace()

	//get correct connector for given system
	switch r.system {
	case GIT:
		r.connector = &GitConnector{}
	}
	r.connector.Load(r.remote, r.local)
}

func (r *Repository) AllCommits() map[string]*Commit {
	return r.connector.Commits()
}

func (r *Repository) AllDevelopers() map[string]*Developer {
	return r.connector.Developers()
}

//Lookup for a single commit
func (r *Repository) FindCommit(id string) *Commit {
	//TODO implement lockup without questioning all commits
	commits := r.AllCommits()
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
		if file, ok := commit.files[fileId]; ok {
			log.Printf("found file %s (%s) in commit %s", fileId, file.path, commitId)
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
