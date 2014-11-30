package analyzer

import (
	"fmt"
	"github.com/jochil/vcs"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type LineSumDiff struct {
	NumAdded   uint
	NumRemoved uint
}

func (d *LineSumDiff) Add(d2 LineSumDiff) {
	d.NumAdded += d2.NumAdded
	d.NumRemoved += d2.NumRemoved
}

func (d LineSumDiff) String() string {
	return fmt.Sprintf("Added: %d, Removed %d", d.NumAdded, d.NumRemoved)
}

//calculates line diff for a developer
func CalcLineDiff(dev *vcs.Developer) LineSumDiff {
	lineSumDiff := LineSumDiff{0, 0}
	for _, commit := range dev.Commits {
		lineSumDiff.Add(CalcLineDiffCommit(commit))
	}
	return lineSumDiff
}

//calculates line diff for a commit
func CalcLineDiffCommit(commit *vcs.Commit) LineSumDiff {

	lineSumDiff := LineSumDiff{0, 0}

	//this is the first commit in the vcs, so count all lines added
	if len(commit.Parents) == 0 {
		for path, file := range commit.Files {
			if Filter.ValidExtension(path) {
				lineSumDiff.Add(execLineDiff("", file.Content()))
			}
		}
	} else {

		//TODO handle merge
		if len(commit.Parents) > 1 {

		} else {
			//get one and only parent commit
			var parentCommit *vcs.Commit
			for _, parentCommit = range commit.Parents {
				break
			}

			from := ""
			//--what is about deleted files?
			for path, file := range commit.Files {
				if Filter.ValidExtension(path) {
					if parentFile := parentCommit.FileByPath(path); parentFile != nil {
						from = parentFile.Content()
					}
					lineSumDiff.Add(execLineDiff(from, file.Content()))
				}
			}
		}
	}

	return lineSumDiff
}

//executes the concrete line diff based on two string
func execLineDiff(from string, to string) LineSumDiff {

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(from, to, true)

	lineDiff := LineSumDiff{0, 0}

	//todo replace -1,0,1 with constants
	for _, diff := range diffs {
		switch diff.Type {
		case -1:
			lineDiff.NumRemoved++
		case 1:
			lineDiff.NumAdded++
		}
	}

	return lineDiff
}
