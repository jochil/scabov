package vcs

import (
	"fmt"
)

//Developer represents a single developer
type Developer struct {
	ident	string
	name 	string
	commits map[string]*Commit
}

func (d *Developer) String() string {
	return fmt.Sprintf("%s (%s): %d Commits\n", d.name, d.ident, len(d.commits))
}
