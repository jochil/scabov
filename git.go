package vcs

import (
	"github.com/libgit2/git2go"
	"fmt"
	"io"
	"crypto/sha1"
	"os"
	"log"
	"path/filepath"
)

//GitLoadCommits get commits from an external repository
func GitLoadCommits(url string) []*Commit {
	repo := getGitRepo(url)
	commits := readGitCommits(repo)
	log.Printf("loaded %d commits from git repo", len(commits))

	return commits
}

//readGitCommits loads all commits from a given repository and transform them to the
//internal Commit type
func readGitCommits(repo *git.Repository) []*Commit {

	//get object database
	odb, err := repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	//iterate over objects and convert them to Commit type
	commits := []*Commit{}
	err = odb.ForEach(func(oid *git.Oid) error {
		obj, err := repo.Lookup(oid)
		if err != nil {
			log.Fatal(err)
		}

		//filter commits
		switch obj := obj.(type) {
		case *git.Commit:
			author := obj.Author()
			dev := &Developer{author.Email} //author.Name
			commit := &Commit{dev: dev, message: obj.Message(), date: author.When}
			commits = append(commits, commit)

		}
		return nil
	})

	return commits
}

func getGitRepo(url string) *git.Repository {

	//get hash from repo url
	h := sha1.New()
	io.WriteString(h, url)
	dir := fmt.Sprintf("%x", h.Sum(nil))

	//get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(cwd)
	}

	//TODO ensure file/dir handling
	path := filepath.Join(cwd, "workspace", dir)

	repo, err := git.OpenRepository(path)
	if err != nil {
		os.RemoveAll(path)
		return cloneGitRepo(url, path)
	} else {
		log.Printf("opened git repo in %s", repo.Path())
		return repo
	}
}

func cloneGitRepo(url string, path string) *git.Repository{
	checkoutOpts := &git.CheckoutOpts{Strategy: git.CheckoutForce}
	cloneOpts := &git.CloneOptions{CheckoutOpts: checkoutOpts}
	repo, err := git.Clone(url, path, cloneOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cloned git repo to %s", repo.Path())
	return repo
}
