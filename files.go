package analyzer

import (
	"fmt"
	"github.com/jochil/vcs"
)

type FileSumDiff struct {
	NumAdded   int
	NumRemoved int
	NumMoved   int
	NumChanged int
}

func (d *FileSumDiff) Add(d2 FileSumDiff) {
	d.NumAdded += d2.NumAdded
	d.NumRemoved += d2.NumRemoved
	d.NumMoved += d2.NumMoved
	d.NumChanged += d2.NumChanged
}

func (d FileSumDiff) String() string {
	return fmt.Sprintf("Added: %d, Removed %d, Moved: %d, Changed: %d",
		d.NumAdded, d.NumRemoved, d.NumMoved, d.NumChanged)
}

func CalcFileDiff(dev *vcs.Developer) FileSumDiff {
	fileSumDiff := FileSumDiff{0, 0, 0, 0}
	for _, commit := range dev.Commits {
		fileSumDiff.Add(CalcFileDiffCommit(commit))
	}
	return fileSumDiff
}

func CalcFileDiffCommit(commit *vcs.Commit) FileSumDiff {

	return FileSumDiff{
		NumAdded:   len(commit.AddedFiles),
		NumRemoved: len(commit.RemovedFiles),
		NumMoved:   len(commit.MovedFiles),
		NumChanged: len(commit.ChangedFiles),
	}
}
