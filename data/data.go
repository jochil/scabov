package data

import (
	"github.com/jochil/analyzer"
	"fmt"
)

type DevData struct {
	Dev              string //TODO replace with pointer to vcs.Dev?
	Commits          uint16
	LineDiff         analyzer.LineSumDiff
}

func (d *DevData) String() string {
	return fmt.Sprintf(
		"%s: \n\t Commits: %d\n\t LineDiff: %s",
		d.Dev,
		d.Commits,
		d.LineDiff,
	)
}
