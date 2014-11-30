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

	Parents  map[string]*Commit
	Children map[string]*Commit
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s by %s: %q @ %v\n-> Removed: %d, Changed: %d, Added: %d, Renamed: %d",
		c.Id, c.Developer, c.Message, c.Date,
		len(c.RemovedFiles), len(c.ChangedFiles), len(c.AddedFiles), len(c.MovedFiles),
	)
}

func (c *Commit) FileByPath(path string) *File {
	if file, exists := c.Files[path]; exists {
		return file
	}
	return nil
}

func (c *Commit) FileById(id string) *File {
	for _, file := range c.Files {
		if file.Id == id {
			return file
		}
	}
	return nil
}
