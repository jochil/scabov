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
	return fmt.Sprintf("%s: %q [%v]\n", c.dev, c.message, c.date)
}
