package vcs

import (
	"time"
	"fmt"
)

//Commit represents a single commit
type Commit struct {
	Id      	string
	Date    	time.Time
	Message 	string
	Developer 	*Developer
	Files 		map[string]*File
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s by %s: %q @ %v", c.Id, c.Developer, c.Message, c.Date)
}
