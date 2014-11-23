package vcs

import (
	git "github.com/libgit2/git2go"
	"os"
	"log"
	"path"
)

//common interface for all vcs connectors
type Connector interface {
	Load(remote string, local string) (map[string]*Commit, map[string]*Developer)
}

//internal struct for this connecotr
type GitConnector struct {
	repo       *git.Repository
	localPath   string
	storagePath string
	commits    map[string]*Commit
	developers map[string]*Developer
	files      map[string]*File
}

//loads an existing repository or clone it from remote
func (c *GitConnector) Load(remote string, local string) (map[string]*Commit, map[string]*Developer) {

	repo, err := git.OpenRepository(local)

	if err != nil {
		os.RemoveAll(local)
		c.repo = c.cloneGitRepo(remote, local)
		log.Printf("cloned git repo from %s to %s", remote, local)
	} else {
		log.Printf("opened git repo in %s", local)
		c.repo = repo
	}

	c.localPath = local
	c.storagePath = path.Join(c.localPath, "file_storage")
	var fm os.FileMode = 0700
	if err := os.MkdirAll(c.storagePath, fm); err != nil {
		log.Fatalln("unable to create storage directory %s", c.storagePath)
	}

	c.developers = map[string]*Developer{}
	c.commits = map[string]*Commit{}
	c.files = map[string]*File{}

	c.fetchAll()

	log.Printf("loaded %d commits, %d delevopers and %d files from git repo",
		len(c.commits), len(c.developers), len(c.files))

	return c.commits, c.developers
}

//return all commits
func (c *GitConnector) fetchAll() {

	if c.repo == nil {
		log.Fatal("no git repository loaded")
	}

	//get object database
	odb, err := c.repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	err = odb.ForEach((git.OdbForEachCallback)(func(oid *git.Oid) error {
		obj, err := c.repo.Lookup(oid)
		if err != nil {
			log.Fatal(err)
		}

		//filter commits
		switch obj := obj.(type) {
		case *git.Commit:
			if _, exists := c.commits[obj.Id().String()]; exists == false {
				c.createCommit(obj)
			}
		}
		return nil
	}))
}

/*
creates an internal commit object based on the git2go commit
recursively create objects for parent commits
 */
func (c GitConnector) createCommit(gitCommit *git.Commit) *Commit {

	author := gitCommit.Author()
	dev, exists := c.developers[author.Email]
	if !exists {
		dev = &Developer{ Id: author.Email, Email: author.Email, Name: author.Name, Commits: map[string]*Commit{}}
		c.developers[author.Email] = dev
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

	//iterate over files
	tree, _ := gitCommit.Tree()
	tree.Walk((git.TreeWalkCallback)(func(filepath string, entry *git.TreeEntry) int {
		if entry.Type == git.ObjectBlob {
			fileId := entry.Id.String()
			blob, err := c.repo.LookupBlob(entry.Id)
			if err != nil {
				log.Fatalf("unable to lookup file %s", filepath+entry.Name)
			} else {
				if file, exists := c.files[fileId]; exists {
					commit.Files[fileId] = file
				}else if Filter.ValidExtension(entry.Name) {
					fileStorage := path.Join(c.storagePath, fileId)

					storeFile(fileStorage, blob.Contents())

					file := &File{Id: fileId, Path: filepath + entry.Name, Size: blob.Size(), StoragePath: fileStorage}
					commit.Files[fileId] = file
					c.files[fileId] = file
				}
			}
		}
		return 0
	}))

	c.commits[gitCommit.Id().String()] = commit
	dev.Commits[gitCommit.Id().String()] = commit

	//iterate over parent commits and create or reference them
	for n := uint(0); n < gitCommit.ParentCount(); n++ {
		parentGitCommit := gitCommit.Parent(n)
		var parentCommit *Commit;
		if parentCommit, exists = c.commits[parentGitCommit.Id().String()]; !exists {
			parentCommit = c.createCommit(parentGitCommit)
		}
		commit.Parents[parentGitCommit.Id().String()] = parentCommit
		parentCommit.Children[commit.Id] = commit

	}

	return commit

}

func (c GitConnector) cloneGitRepo(remote string, local string) *git.Repository {
	checkoutOpts := &git.CheckoutOpts{Strategy:git.CheckoutForce}
	cloneOpts := &git.CloneOptions{CheckoutOpts: checkoutOpts, Bare: true}
	repo, err := git.Clone(remote, local, cloneOpts)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
