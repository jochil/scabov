package vcs

import (
	"time"
	"fmt"
)

//Commit represents a single commit
type Commit struct {
	id       string
	date     time.Time
	message  string
	dev      *Developer
	files    map[string]*File
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s by %s: %q @ %v", c.id, c.dev, c.message, c.date)
}
