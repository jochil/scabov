package vcs

type LineDiff struct {
	Added   int
	Removed int
}

func (diff *LineDiff) Add(diff2 LineDiff) {
	diff.Added += diff2.Added
	diff.Removed += diff2.Removed
}
