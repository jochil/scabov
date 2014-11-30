package analyzer

import (
	"fmt"
	"github.com/jochil/vcs"
	"github.com/sergi/go-diff/diffmatchpatch"
	"strings"
)

type LineSumDiff struct {
	NumAdded   int
	NumRemoved int
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
		for _, file := range commit.Files {
			lineSumDiff.Add(execLineDiff("", file.Content()))
		}
	} else {

		//TODO handle merge
		if len(commit.Parents) > 1 {

		} else {
			for _, file := range commit.AddedFiles {
				lineSumDiff.Add(execLineDiff("", file.Content()))
			}
			for _, file := range commit.RemovedFiles {
				lineSumDiff.Add(execLineDiff(file.Content(), ""))
			}
			for _, file := range commit.ChangedFiles {
				for _, parentFile := range file.Parents {
					lineSumDiff.Add(execLineDiff(parentFile.Content(), file.Content()))
					//TODO how to handle multiple parent files?
					break
				}
			}
		}
	}

	return lineSumDiff
}

//executes the concrete line diff based on two strings
func execLineDiff(from string, to string) LineSumDiff {

	dmp := diffmatchpatch.New()

	//use line mode
	//see: https://code.google.com/p/google-diff-match-patch/wiki/LineOrWordDiffs
	lineText1, lineText2, lineArray := dmp.DiffLinesToChars(from, to)
	diffs := dmp.DiffMain(lineText1, lineText2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiff := LineSumDiff{0, 0}

	for _, diff := range diffs {
		lines := strings.Count(diff.Text, "\n")

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			lineDiff.NumRemoved += lines
		case diffmatchpatch.DiffInsert:
			lineDiff.NumAdded += lines
		}
	}

	return lineDiff
}
