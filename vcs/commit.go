package vcs

import (
	"fmt"
	"time"
)

//Commit represents a single commit
type Commit struct {
	Id        string
	Date      time.Time
	Message   string
	Developer *Developer

	Files        map[string]*File
	RemovedFiles map[string]*File
	ChangedFiles map[string]*File
	AddedFiles   map[string]*File
	MovedFiles   map[string]*File

	LineDiff LineDiff

	Parents  map[string]*Commit
	Children map[string]*Commit
}

func NewCommit(id string, message string, date time.Time, dev *Developer) *Commit {
	return &Commit{
		Id:           id,
		Developer:    dev,
		Message:      message,
		Date:         date,
		Files:        map[string]*File{},
		ChangedFiles: map[string]*File{},
		RemovedFiles: map[string]*File{},
		AddedFiles:   map[string]*File{},
		MovedFiles:   map[string]*File{},
		Parents:      map[string]*Commit{},
		LineDiff:     LineDiff{0, 0},
		Children:     map[string]*Commit{},
	}
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s by %s: %q @ %v\n-> Removed: %d, Changed: %d, Added: %d, Renamed: %d",
		c.Id, c.Developer, c.Message, c.Date,
		len(c.RemovedFiles), len(c.ChangedFiles), len(c.AddedFiles), len(c.MovedFiles),
	)
}
