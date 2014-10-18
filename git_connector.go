package vcs

import (
	"github.com/libgit2/git2go"
	"os"
	"log"
)

//internal struct for this connecotr
type GitConnector struct {
	Repo *git.Repository
	Commits map[string]*Commit
	Devs map[string]*Developer
}

//loads an existing repository or clone it from remote
func (c *GitConnector) Load(remote string, local string) {

	//TODO init values, should be done in constructor
	c.Devs = map[string]*Developer{}

	repo, err := git.OpenRepository(local)
	if err != nil {
		os.RemoveAll(local)
		c.Repo = c.cloneGitRepo(remote, local)
	} else {
		log.Printf("opened git repo in %s", repo.Path())
		c.Repo = repo
	}
}

func (c GitConnector) AllDevelopers() map[string]*Developer {
	//TODO replace this with a more efficient way
	if c.Commits == nil {
		c.AllCommits()
	}

	log.Printf("found %d developers in git repo", len(c.Devs))
	return c.Devs
}

//return all commits
func (c *GitConnector) AllCommits() map[string]*Commit {

	if c.Repo == nil {
		log.Fatal("no git repository loaded")
	}

	if c.Commits != nil {
		return c.Commits
	}

	//get object database
	odb, err := c.Repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	//iterate over objects and convert them to Commit type
	commits := map[string]*Commit{}
	err = odb.ForEach(func(oid *git.Oid) error {
		obj, err := c.Repo.Lookup(oid)
		if err != nil {
			log.Fatal(err)
		}

		//filter commits
		switch obj := obj.(type) {
		case *git.Commit:
			author := obj.Author()

			dev, exists := c.Devs[author.Email]
			if !exists {
				dev = &Developer{ Id: author.Email, Email: author.Email, Name: author.Name, Commits: map[string]*Commit{}}
				c.Devs[author.Email] = dev
			}

			commit := &Commit{Id: obj.Id().String(), Developer: dev, Message: obj.Message(), Date: author.When, Files: map[string]*File{}}
			log.Println(commit)

			//iterate over files
			tree, _ := obj.Tree()
			tree.Walk(func(path string, entry *git.TreeEntry) int{
				if entry.Type == git.ObjectBlob {
					fileId := entry.Id.String()
					blob, err := c.Repo.LookupBlob(entry.Id)
					if err != nil {
						log.Fatalf("unable to file %s", path + entry.Name)
					} else {
						commit.Files[fileId] = &File{Id: fileId, Path: path + entry.Name, Size: blob.Size(), Contents: blob.Contents()}
						log.Println("\t", path + entry.Name, entry.Id)
					}

				}
				return 0
			})

			commits[obj.Id().String()] = commit
			dev.Commits[obj.Id().String()] = commit

		}
		return nil
	})
	log.Printf("loaded %d commits from git repo", len(commits))
	c.Commits = commits //cache commits
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
