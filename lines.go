package analyzer

import (
	"github.com/jochil/vcs"
	"github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
)

type LineSumDiff struct {
	Added     uint32
	Removed   uint32
	Untouched uint32
}

func (d *LineSumDiff) Add(d2 LineSumDiff) {
	d.Added += d2.Added
	d.Removed += d2.Removed
	d.Untouched += d2.Untouched
}

func (d LineSumDiff) String() string {
	return fmt.Sprintf("Added: %d, Removed %d, Untouched: %d", d.Added, d.Removed, d.Untouched)
}

//calculates line diff for a developer
func CalcLineDiff(dev *vcs.Developer) LineSumDiff {
	lineSumDiff := LineSumDiff{0, 0, 0}
	for _, commit := range dev.Commits {
		lineSumDiff.Add(CalcLineDiffCommit(commit))
	}
	return lineSumDiff
}

//calculates line diff for a commit
func CalcLineDiffCommit(commit *vcs.Commit) LineSumDiff {

	lineSumDiff := LineSumDiff{0, 0, 0}

	//this is the first commit in the vcs, so count all lines added
	if len(commit.Parents) == 0 {
		for _, file := range commit.Files {
			if Filter.ValidExtension(file.Path) {
				lineSumDiff.Add(execLineDiff("", file.Content()))
			}
		}
	}  else {

		//TODO handle merge
		if len(commit.Parents) > 1 {

		} else {
			//get one and only parent commit
			var parentCommit *vcs.Commit
			for _, parentCommit = range commit.Parents {
				break;
			}

			from := ""
			//--what is about deleted files?
			for _, file := range commit.Files {
				if Filter.ValidExtension(file.Path) {
					if parentFile := parentCommit.FileByPath(file.Path); parentFile != nil {
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

	lineDiff := LineSumDiff{0, 0, 0}

	//todo replace -1,0,1 with constants
	for _, diff := range diffs {
		switch diff.Type {
		case -1:
			lineDiff.Removed++
		case 1:
			lineDiff.Added++
		case 0:
			lineDiff.Untouched++
		}
	}

	//if added and removed == 0 the file was not edited in this commit,
	if lineDiff.Added == 0 && lineDiff.Removed == 0 {
		lineDiff.Untouched = 0
	}

	return lineDiff
}
