package vcs

import (
	"fmt"
)

//Developer represents a single developer
type Developer struct {
	Id		string
	Name 	string
	Email	string
	Commits map[string]*Commit
}

func (d *Developer) String() string {
	return fmt.Sprintf("%s (%s)", d.Name, d.Email)
}
