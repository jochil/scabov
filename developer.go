package vcs

import (
	"fmt"
)

//Developer represents a single developer
type Developer struct {
	id		string
	name 	string
	email	string
	commits map[string]*Commit
}

func (d *Developer) String() string {
	return fmt.Sprintf("%s (%s)", d.name, d.email)
}
