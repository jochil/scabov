package analyzer

import (
	"fmt"
	"github.com/jochil/scabov/vcs"
)

type FunctionHistory struct {
	Name             string
	File             string
	lifetime         int
	changes          int
	removed          bool
	latestHash       string
	firstComplexity  int
	latestComplexity int
	firstSize        int
	latestSize       int
}

func NewFunctionHistory(function Function, file string) *FunctionHistory {

	cyclo := CyclomaticComplexity(function.CFG)
	return &FunctionHistory{
		Name:             function.Name,
		File:             file,
		changes:          0,
		lifetime:         0,
		removed:          false,
		latestHash:       function.Hash,
		firstComplexity:  cyclo,
		latestComplexity: cyclo,
		firstSize:        function.NumNodes,
		latestSize:       function.NumNodes,
	}
}

func (history *FunctionHistory) Remove() {
	history.lifetime++
	history.changes++
	history.removed = true
}

func (history *FunctionHistory) Change(function Function) {

	if history.removed == false {
		if history.latestHash != function.Hash {
			history.latestHash = function.Hash
			history.latestSize = function.NumNodes
			history.latestComplexity = CyclomaticComplexity(function.CFG)
			history.changes++
		}
		history.lifetime++
	}
}

func (history *FunctionHistory) Beat() {
	if history.removed == false {
		history.lifetime++
	}
}

func (history *FunctionHistory) Growth() (size float64, complexity float64) {
	if history.lifetime == 0 {
		return 0.0, 0.0
	}

	size = float64(history.latestSize-history.firstSize) / float64(history.lifetime)
	complexity = float64(history.latestComplexity-history.firstComplexity) / float64(history.lifetime)

	return size, complexity
}

func (history *FunctionHistory) Stability() float64 {
	if history.lifetime == 0 {
		return 0
	}
	return float64(history.lifetime-history.changes) / float64(history.lifetime)
}

func (history *FunctionHistory) String() string {
	sizeGrowthh, complexityGrowth := history.Growth()

	return fmt.Sprintf("%s L=%d, C=%d, S=%.4f, GS=%.4f, GC=%.4f",
		history.Name, history.lifetime, history.changes, history.Stability(),
		sizeGrowthh, complexityGrowth)
}

type FileHistory map[string]*FunctionHistory

func (fileHistory FileHistory) Beat() {
	for _, functionHistory := range fileHistory {
		functionHistory.Beat()
	}
}

func (fileHistory FileHistory) Remove() {
	for _, functionHistory := range fileHistory {
		functionHistory.Remove()
	}
}

func (fileHistory FileHistory) Rename(filename string) {
	for _, functionHistory := range fileHistory {
		functionHistory.File = filename
	}
}

var History = map[string]FileHistory{}

func LoadHistory(repo *vcs.Repository) {
	readCommit(repo.FirstCommit())
}

func readCommit(commit *vcs.Commit) {
	readHistory(commit)
	for _, child := range commit.Children {
		readCommit(child)
		//TODO just a workaround to avoid endless processing... shoud be removed
		break
	}
}

func readHistory(commit *vcs.Commit) {
	parser := NewParser()

	//handle the beat for files that was not part of this commit
	for filename, fileHistory := range History {
		_, listed := commit.Files[filename]
		_, removed := commit.RemovedFiles[filename]
		_, moved := commit.MovedFiles[filename]

		if listed == false && removed == false && moved == false {
			fileHistory.Beat()
		}
	}

	//handle moved files
	for oldFilename, newFilename := range commit.MovedFiles {

		//update history to current filename
		if _, ok := History[newFilename]; ok == false {
			fileHistory := History[oldFilename]
			History[newFilename] = fileHistory
			fileHistory.Rename(newFilename)
			delete(History, oldFilename)

			//file was just moved not changed, so handle the beat
			if _, ok := commit.ChangedFiles[newFilename]; ok == false {
				fileHistory.Beat()
			}
		}
	}

	//find modified functions
	for filename, file := range commit.ChangedFiles {

		if _, ok := History[filename]; ok == false {
			History[filename] = FileHistory{}
		}

		fileHistory := History[filename]
		functions := parser.Functions(file)

		//search for (un)changed function
		for name, function := range functions {
			if history, ok := fileHistory[name]; ok {
				history.Change(function)
			} else {
				fileHistory[name] = NewFunctionHistory(function, filename)
			}
		}

		//search removed functions
		for name, history := range fileHistory {
			if _, ok := functions[name]; ok == false {
				history.Remove()
			}
		}
	}

	//handle removed files
	for filename, _ := range commit.RemovedFiles {
		fileHistory := History[filename]
		fileHistory.Remove()
	}

	//find new functions
	for filename, file := range commit.AddedFiles {

		fileHistory := FileHistory{}
		History[filename] = fileHistory

		functions := parser.Functions(file)
		for name, function := range functions {
			fileHistory[name] = NewFunctionHistory(function, filename)
		}
	}
}
