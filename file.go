package vcs

import (
	"io/ioutil"
	"os"
	"log"
)

// stores file to local path
func storeFile(path string, content []byte) {
	var perm os.FileMode = 0600
	if err := ioutil.WriteFile(path, content, perm); err != nil {
		log.Fatalln(err)
	}
}

type File struct {
	Id          string
	Path        string
	Size        int64
	StoragePath string
}

// returns file content as string
func (f *File) Content() string {

	if content, err := ioutil.ReadFile(f.StoragePath); err != nil {
		panic(err)
	} else {
		return string(content[:])
	}
}


