package analyzer

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func CountLineDiff(from string, to string) (insert uint16, delete uint16, equal uint16) {

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(from, to, true)

	//todo replace -1,0,1 with constants
	for _, diff := range diffs {
		switch diff.Type {
		case -1:
			delete++
		case 1:
			insert++
		case 0:
			equal++
		}
	}

	return insert, delete, equal
}
