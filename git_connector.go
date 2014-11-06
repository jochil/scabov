package vcs

import (
	git "github.com/libgit2/git2go"
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
	c.Commits = map[string]*Commit{}

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
	if len(c.Commits) == 0 {
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

	if len(c.Commits) > 0 {
		return c.Commits
	}

	//get object database
	odb, err := c.Repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	err = odb.ForEach((git.OdbForEachCallback)(func(oid *git.Oid) error {
		obj, err := c.Repo.Lookup(oid)
		if err != nil {
			log.Fatal(err)
		}

		//filter commits
		switch obj := obj.(type) {
		case *git.Commit:
			if _, exists := c.Commits[obj.Id().String()]; exists == false {
				c.createCommit(obj)
			}
		}
		return nil
	}))
	log.Printf("loaded %d commits from git repo", len(c.Commits))
	return c.Commits
}

/*
creates an internal commit object based on the git2go commit
recursively create objects for parent commits
 */
func (c GitConnector) createCommit(gitCommit *git.Commit) *Commit {

	author := gitCommit.Author()
	dev, exists := c.Devs[author.Email]
	if !exists {
		dev = &Developer{ Id: author.Email, Email: author.Email, Name: author.Name, Commits: map[string]*Commit{}}
		c.Devs[author.Email] = dev
	}

	commit := &Commit{
		Id: gitCommit.Id().String(),
		Developer: dev,
		Message: gitCommit.Message(),
		Date: author.When,
		Files: map[string]*File{},
		Parents: map[string]*Commit{},
		Children: map[string]*Commit{},
	}
	log.Println("found commit:", commit)

	//iterate over files
	tree, _ := gitCommit.Tree()
	tree.Walk((git.TreeWalkCallback)(func(path string, entry *git.TreeEntry) int {
		if entry.Type == git.ObjectBlob {
			fileId := entry.Id.String()
			blob, err := c.Repo.LookupBlob(entry.Id)
			if err != nil {
				log.Fatalf("unable to file %s", path+entry.Name)
			} else {
				commit.Files[fileId] = &File{Id: fileId, Path: path+entry.Name, Size: blob.Size(), Contents: blob.Contents()}
			}

		}
		return 0
	}))

	c.Commits[gitCommit.Id().String()] = commit
	dev.Commits[gitCommit.Id().String()] = commit

	//iterate over parent commits and create or reference them
	for n := uint(0); n < gitCommit.ParentCount(); n++ {
		parentGitCommit := gitCommit.Parent(n)
		var parentCommit *Commit;
		if parentCommit, exists = c.Commits[parentGitCommit.Id().String()]; !exists {
			parentCommit = c.createCommit(parentGitCommit)
		}
		commit.Parents[parentGitCommit.Id().String()] = parentCommit
		parentCommit.Children[commit.Id] = commit

	}

	return commit

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
