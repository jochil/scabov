package vcs

import (
	"time"
	"fmt"
)

//Commit represents a single commit
type Commit struct {
	id         string
	date       time.Time
	message    string
	dev        *Developer
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s: %q\n[%v]", c.dev, c.message, c.date)
}

//Developer represents a single developer
type Developer struct {
	name string
}

func (d *Developer) String() string {
	return fmt.Sprintf("%s", d.name)
}
