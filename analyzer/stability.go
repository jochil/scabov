package analyzer

import (
	"fmt"
	"github.com/jochil/scabov/vcs"
	"log"
)

type FunctionHistory struct {
	Name       string
	File       string
	lifetime   int
	changes    int
	removed    bool
	latestHash string
}

func NewFunctionHistory(name string, file string, firstHash string) *FunctionHistory {

	return &FunctionHistory{
		Name:       name,
		File:       file,
		changes:    0,
		lifetime:   0,
		removed:    false,
		latestHash: firstHash,
	}
}

func (history *FunctionHistory) Remove() {
	history.lifetime++
	history.changes++
	history.removed = true
}

func (history *FunctionHistory) Change(hash string) {

	if history.removed == false {
		if history.latestHash != hash {
			history.latestHash = hash
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

func (history *FunctionHistory) Stability() float64 {
	if history.lifetime == 0 {
		return 0
	}
	return float64(history.lifetime-history.changes) / float64(history.lifetime)
}

func (history *FunctionHistory) String() string {
	return fmt.Sprintf("%s L=%d, C=%d, S=%.4f", history.Name, history.lifetime, history.changes, history.Stability())
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

var overallHistory = map[string]FileHistory{}

func CalcFunctionStability(repo *vcs.Repository) (float64, map[string]FileHistory) {
	readCommit(repo.FirstCommit())

	count := 0.0
	sum := 0.0
	for _, fileHistory := range overallHistory {
		for _, history := range fileHistory {
			count++
			sum += history.Stability()
		}
	}

	stability := sum / count

	return stability, overallHistory
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
	for filename, fileHistory := range overallHistory {
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
		if _, ok := overallHistory[newFilename]; ok == false {
			fileHistory := overallHistory[oldFilename]
			overallHistory[newFilename] = fileHistory
			fileHistory.Rename(newFilename)
			delete(overallHistory, oldFilename)

			//file was just moved not changed, so handle the beat
			if _, ok := commit.ChangedFiles[newFilename]; ok == false {
				fileHistory.Beat()
			}
		}
	}

	//find modified functions
	for filename, file := range commit.ChangedFiles {

		if _, ok := overallHistory[filename]; ok == false {
			overallHistory[filename] = FileHistory{}
		}

		fileHistory := overallHistory[filename]
		functions := parser.Functions(file)

		//search for (un)changed function
		for name, function := range functions {
			if history, ok := fileHistory[name]; ok {
				history.Change(function.Hash)
			} else {

				if fileHistory == nil {
					log.Println(filename)
					log.Println("FUCK")
				}

				fileHistory[name] = NewFunctionHistory(name, filename, function.Hash)
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
		fileHistory := overallHistory[filename]
		fileHistory.Remove()
	}

	//find new functions
	for filename, file := range commit.AddedFiles {

		fileHistory := FileHistory{}
		overallHistory[filename] = fileHistory

		functions := parser.Functions(file)
		for name, function := range functions {
			fileHistory[name] = NewFunctionHistory(name, filename, function.Hash)
		}
	}
}
