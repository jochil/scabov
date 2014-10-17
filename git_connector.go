package vcs

import (
	"github.com/libgit2/git2go"
	"os"
	"log"
)

//internal struct for this connecotr
type GitConnector struct {
	repo *git.Repository
}

//loads an existing repository or clone it from remote
func (c *GitConnector) Load(remote string, local string) {
	repo, err := git.OpenRepository(local)
	if err != nil {
		os.RemoveAll(local)
		c.repo = c.cloneGitRepo(remote, local)
	} else {
		log.Printf("opened git repo in %s", repo.Path())
		c.repo = repo
	}
}

//return all commits
func (c GitConnector) Commits() map[string]*Commit {

	if c.repo == nil {
		log.Fatal("no git repository loaded")
	}

	//get object database
	odb, err := c.repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	//iterate over objects and convert them to Commit type
	commits := map[string]*Commit{}
	err = odb.ForEach(func(oid *git.Oid) error {
		obj, err := c.repo.Lookup(oid)
		if err != nil {
			log.Fatal(err)
		}

		//filter commits
		switch obj := obj.(type) {
		case *git.Commit:
			author := obj.Author()
			dev := &Developer{author.Email} //author.Name
			commit := &Commit{dev: dev, message: obj.Message(), date: author.When}
			commits[obj.Id().String()] = commit

		}
		return nil
	})
	log.Printf("loaded %d commits from git repo", len(commits))
	return commits
}

func (c GitConnector) cloneGitRepo(remote string, local string) *git.Repository {

	checkoutOpts := &git.CheckoutOpts{Strategy:git.CheckoutForce}
	cloneOpts := &git.CloneOptions{CheckoutOpts: checkoutOpts}
	repo, err := git.Clone(remote, local, cloneOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cloned git repo to %s", repo.Path())
	return repo
}
