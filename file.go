package vcs

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
)

type FileDiff struct {
	Added   int
	Removed int
	Changed int
}

type File struct {
	Id          string
	Size        int64
	StoragePath string
	Parents     []*File
}

func (f *File) String() string {
	return fmt.Sprintf("%s[%d bytes]", f.Id, f.Size)
}

// returns file content as string
func (f *File) Content() string {

	if content, err := ioutil.ReadFile(f.StoragePath); err != nil {
		panic(err)
	} else {
		return string(content[:])
	}
}

// stores file to local path
func storeFile(path string, content []byte) {

	if _, err := os.Stat(path); err != nil {
		var perm os.FileMode = 0600
		if err := ioutil.WriteFile(path, content, perm); err != nil {
			log.Fatalln(err)
		}
	}
}
