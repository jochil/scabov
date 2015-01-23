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

func NewDeveloper(id string, email string, name string) *Developer {
	return &Developer{
		Id:      id,
		Email:   email,
		Name:    name,
		Commits: map[string]*Commit{},
	}
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

func (dev *Developer) ModifiedFiles() []*File {
	files := []*File{}
	for _, commit := range dev.Commits {
		for _, file := range commit.ChangedFiles {
			files = append(files, file)
		}
	}
	return files
}

func (dev *Developer) AddedFiles() []*File {
	files := []*File{}
	for _, commit := range dev.Commits {
		for _, file := range commit.AddedFiles {
			files = append(files, file)
		}
	}
	return files
}

func (dev *Developer) LineDiff() *LineDiff {

	diff := &LineDiff{0, 0}
	for _, commit := range dev.Commits {
		diff.Add(commit.LineDiff)
	}
	return diff
}

func (dev *Developer) FileDiff() *FileDiff {

	diff := &FileDiff{0, 0, 0}
	for _, commit := range dev.Commits {
		diff.Added += len(commit.AddedFiles)
		diff.Removed += len(commit.RemovedFiles)
		diff.Changed += len(commit.ChangedFiles)
		diff.Changed += len(commit.MovedFiles)
	}
	return diff
}

func (dev *Developer) String() string {
	return fmt.Sprintf("%s (%s)", dev.Name, dev.Email)
}
