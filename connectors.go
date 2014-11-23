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

	log.Printf("loaded %d commits, %d delevopers and %d different files from git repo",
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
		dev = &Developer{
			Id: author.Email,
			Email: author.Email,
			Name: author.Name,
			Commits: map[string]*Commit{},
		}
		c.developers[author.Email] = dev
	}

	commit := &Commit{
		Id: gitCommit.Id().String(),
		Developer: dev,
		Message: gitCommit.Message(),
		Date: author.When,
		Files: map[string]*File{},
		ChangedFiles: map[string]*File{},
		RemovedFiles: map[string]*File{},
		NewFiles: map[string]*File{},
		RenamedFiles: map[string]*File{},
		Parents: map[string]*Commit{},
		Children: map[string]*Commit{},
	}

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

	//iterate over files
	tree, _ := gitCommit.Tree()
	tree.Walk((git.TreeWalkCallback)(func(dir string, entry *git.TreeEntry) int {
		if entry.Type == git.ObjectBlob {

			if Filter.ValidExtension(entry.Name) {
				var file *File;
				fileId := entry.Id.String()
				filepath := dir + entry.Name

				if file, exists = c.files[fileId]; exists {
					commit.Files[filepath] = file
				} else {
					file = c.loadFile(entry)
					commit.Files[filepath] = file
					c.files[fileId] = file
				}

				//connect with parent commits
				if len(commit.Parents) == 0 {
					commit.NewFiles[filepath] = file
				} else {
					for _, parentCommit := range commit.Parents {
						if parentFile := parentCommit.FileById(file.Id); parentFile != nil {
							if parentCommit.FileByPath(filepath) == nil {
								commit.RenamedFiles[filepath] = file
							}
						} else if parentFile := parentCommit.FileByPath(filepath); parentFile != nil{
							commit.ChangedFiles[filepath] = file
						} else {
							commit.NewFiles[filepath] = file
						}
					}
				}
				//TODO and deleted files

			}
		}
		return 0
	}))

	return commit

}

func (c GitConnector) loadFile(entry *git.TreeEntry) *File {
	fileId := entry.Id.String()
	var file *File;
	if blob, err := c.repo.LookupBlob(entry.Id); err == nil {
		fileStorage := path.Join(c.storagePath, fileId)
		storeFile(fileStorage, blob.Contents())
		file = &File{Id: fileId, Size: blob.Size(), StoragePath: fileStorage}
	} else {
		log.Fatalf("unable to lookup file %s", entry.Name)
	}
	return file
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
