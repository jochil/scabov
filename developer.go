package vcs

import (
	"fmt"
)

//Developer represents a single developer
type Developer struct {
	Id      string
	Name    string
	Email   string
	Commits map[string]*Commit
}

func (dev *Developer) FirstCommit() *Commit {
	var firstCommit *Commit
	for _, commit := range dev.Commits {
		if firstCommit == nil {
			firstCommit = commit
		}
		if commit.Date.Before(firstCommit.Date) {
			firstCommit = commit
		}
	}

	return firstCommit
}

func (dev *Developer) String() string {
	return fmt.Sprintf("%s (%s)", dev.Name, dev.Email)
}
