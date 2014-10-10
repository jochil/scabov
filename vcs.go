/*
Package vcs
 */
package vcs

import (
	"time"
	"fmt"
)

//Commit represents a single commit
type Commit struct {
	date         time.Time
	message      string
	developer    *Developer
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s: %q\n[%v]", c.developer, c.message, c.date)
}

//Developer represents a single developer
type Developer struct {
	name string
}

func (d *Developer) String() string {
	return fmt.Sprintf("%s", d.name)
}

// SampleCommit returns a single commit
func SampleCommit() *Commit {
	dev := &Developer{"Jochen"}
	c := &Commit{}
	c.developer = dev
	c.message = "sample Commit"
	c.date = time.Now()
	return c
}
