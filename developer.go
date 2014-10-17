package vcs

import (
	"fmt"
)

//Developer represents a single developer
type Developer struct {
	name string
}

func (d *Developer) String() string {
	return fmt.Sprintf("%s", d.name)
}
