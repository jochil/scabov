package vcs

import (
	"time"
	"fmt"
)

//Commit represents a single commit
type Commit struct {
	Id          string
	Date        time.Time
	Message     string
	Developer   *Developer
	Files       map[string]*File
	Parents     map[string]*Commit
	Children    map[string]*Commit
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s by %s: %q @ %v", c.Id, c.Developer, c.Message, c.Date)
}

func (c *Commit) FileByPath(path string) *File {
	for _, file := range c.Files {
		if file.Path == path {
			return file
		}
	}
	return nil
}
