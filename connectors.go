package vcs

import (
	git "github.com/libgit2/git2go"
	"log"
	"os"
	"path"
)

//common interface for all vcs connectors
type Connector interface {
	Load(remote string, local string) (map[string]*Commit, map[string]*Developer)
}

//internal struct for this connecotr
type GitConnector struct {
	repo        *git.Repository
	localPath   string
	storagePath string
	commits     map[string]*Commit
	developers  map[string]*Developer
	files       map[string]*File
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
			Id:      author.Email,
			Email:   author.Email,
			Name:    author.Name,
			Commits: map[string]*Commit{},
		}
		c.developers[author.Email] = dev
	}

	commit := &Commit{
		Id:           gitCommit.Id().String(),
		Developer:    dev,
		Message:      gitCommit.Message(),
		Date:         author.When,
		Files:        map[string]*File{},
		ChangedFiles: map[string]*File{},
		RemovedFiles: map[string]*File{},
		AddedFiles:   map[string]*File{},
		MovedFiles:   map[string]*File{},
		Parents:      map[string]*Commit{},
		Children:     map[string]*Commit{},
	}

	c.commits[gitCommit.Id().String()] = commit
	dev.Commits[gitCommit.Id().String()] = commit

	//iterate over parent commits and create or reference them
	for n := uint(0); n < gitCommit.ParentCount(); n++ {
		parentGitCommit := gitCommit.Parent(n)
		var parentCommit *Commit
		if parentCommit, exists = c.commits[parentGitCommit.Id().String()]; !exists {
			parentCommit = c.createCommit(parentGitCommit)
		}
		commit.Parents[parentGitCommit.Id().String()] = parentCommit
		parentCommit.Children[commit.Id] = commit

	}

	tree, _ := gitCommit.Tree()

	if gitCommit.ParentCount() == 0 {
		c.loadTreeDiffToCommit(commit, &git.Tree{}, tree)
	}

	for n := uint(0); n < gitCommit.ParentCount(); n++ {
		parentGitCommit := gitCommit.Parent(n)
		parentTree, _ := parentGitCommit.Tree()
		c.loadTreeDiffToCommit(commit, parentTree, tree)
	}
	return commit
}

func (c GitConnector) loadTreeDiffToCommit(commit *Commit, parentTree *git.Tree, newTree *git.Tree) {
	diffOpt, err := git.DefaultDiffOptions()
	if err != nil {
		log.Fatalln(err)
	}

	diff, err := c.repo.DiffTreeToTree(parentTree, newTree, &diffOpt)
	if err != nil {
		log.Fatalln(err)
	}

	findOpts, _ := git.DefaultDiffFindOptions()
	findOpts.Flags += git.DiffFindAll
	diff.FindSimilar(&findOpts)
	err = diff.ForEach((git.DiffForEachFileCallback)(func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		if Filter.ValidExtension(delta.NewFile.Path) {

			var file *File
			fileId := delta.NewFile.Oid.String()
			filepath := delta.NewFile.Path

			var oldFile *File
			oldFileId := delta.OldFile.Oid.String()
			oldFilepath := delta.OldFile.Path

			var exists bool

			if oldFile, exists = c.files[oldFileId]; exists == false && delta.OldFile.Oid.IsZero() == false {
				oldFile = c.loadFile(delta.OldFile.Oid)
				c.files[oldFileId] = oldFile
			}

			if file, exists = c.files[fileId]; exists {
				commit.Files[fileId] = file
			} else if delta.NewFile.Oid.IsZero() {
				commit.RemovedFiles[oldFilepath] = oldFile
			} else {
				file = c.loadFile(delta.NewFile.Oid)
				commit.Files[filepath] = file
				c.files[fileId] = file
			}

			var status string
			switch delta.Status {

			case git.DeltaModified:
				status = "Modified"
				commit.ChangedFiles[filepath] = file
				file.Parents = append(file.Parents, oldFile)
			case git.DeltaAdded:
				status = "Added"
				commit.AddedFiles[filepath] = file
			case git.DeltaRenamed:
				status = "Renamed"
				commit.MovedFiles[filepath] = file
				if delta.Similarity != 100 {
					status += "/Modified"
					commit.ChangedFiles[filepath] = file
					file.Parents = append(file.Parents, oldFile)
				}
			case git.DeltaDeleted:
				status = "Deleted"
			case git.DeltaCopied:
				status = "Copied"

			default:
				status = "unknown"
			}

			//log.Println(status, delta.Similarity, delta.OldFile.Path, delta.OldFile.Oid, delta.NewFile.Path, delta.NewFile.Oid)
		}

		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				return nil
			}, nil
		}, nil
	}), git.DiffDetailLines)
}

func (c GitConnector) loadFile(oid *git.Oid) *File {
	var file *File
	if blob, err := c.repo.LookupBlob(oid); err == nil {
		fileStorage := path.Join(c.storagePath, oid.String())
		storeFile(fileStorage, blob.Contents())
		file = &File{Id: oid.String(), Size: blob.Size(), StoragePath: fileStorage}
	} else {
		log.Fatalf("unable to lookup file %s", oid)
	}
	return file
}

func (c GitConnector) cloneGitRepo(remote string, local string) *git.Repository {
	checkoutOpts := &git.CheckoutOpts{Strategy: git.CheckoutForce}
	cloneOpts := &git.CloneOptions{CheckoutOpts: checkoutOpts, Bare: true}
	repo, err := git.Clone(remote, local, cloneOpts)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
