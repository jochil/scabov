package vcs

import (
	"fmt"
	git "github.com/libgit2/git2go"
	"log"
	"os"
	"path"
)

//common interface for all vcs connectors
type Connector interface {
	LoadRemote(path string, workspace string) error
	LoadLocal(path string, workspace string) error
	Developers() map[string]*Developer
	Commits() map[string]*Commit
}

//internal struct for this connecotr
type GitConnector struct {
	repo        *git.Repository
	storagePath string
	commits     map[string]*Commit
	developers  map[string]*Developer
	files       map[string]*File
}

func (c *GitConnector) LoadLocal(path string, workspace string) error {

	repo, err := git.OpenRepository(path)
	if err != nil {
		return err
	}

	if err := c.init(repo, workspace); err != nil {
		return err
	}

	log.Printf("opened local git repo in %s", path)
	return nil
}

//loads an existing repository or clone it from external source
func (c *GitConnector) LoadRemote(path string, workspace string) error {

	//clear workspace
	os.RemoveAll(workspace)
	repo, err := c.cloneGitRepo(path, workspace)
	if err != nil {
		return err
	}

	if err := c.init(repo, workspace); err != nil {
		return err
	}

	log.Printf("cloned remote git repo from %s to %s", path, workspace)

	return nil
}

func (c *GitConnector) init(repo *git.Repository, workspace string) error {

	c.repo = repo
	c.storagePath = workspace

	var fm os.FileMode = 0700
	if err := os.MkdirAll(c.storagePath, fm); err != nil {
		return fmt.Errorf("unable to create storage directory %s", c.storagePath)
	}

	c.developers = map[string]*Developer{}
	c.commits = map[string]*Commit{}
	c.files = map[string]*File{}

	c.fetchAll()

	log.Printf("loaded %d commits, %d delevopers and %d different files from git repo",
		len(c.commits), len(c.developers), len(c.files))

	return nil
}

func (c *GitConnector) Developers() map[string]*Developer {
	return c.developers
}

func (c *GitConnector) Commits() map[string]*Commit {
	return c.commits
}

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
		dev = NewDeveloper(author.Email, author.Email, author.Name)
		c.developers[author.Email] = dev
	}

	commit := NewCommit(gitCommit.Id().String(), gitCommit.Message(), author.When, dev)

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
				commit.MovedFiles[oldFilepath] = filepath
				commit.Files[filepath] = file

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
				if Filter.ValidExtension(delta.NewFile.Path) {
					if line.Origin == git.DiffLineAddition {
						commit.LineDiff.Added++
					} else if line.Origin == git.DiffLineDeletion {
						commit.LineDiff.Removed++
					}
				}

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

func (c GitConnector) cloneGitRepo(external string, local string) (*git.Repository, error) {
	checkoutOpts := &git.CheckoutOpts{Strategy: git.CheckoutForce}
	cloneOpts := &git.CloneOptions{CheckoutOpts: checkoutOpts, Bare: true}
	repo, err := git.Clone(external, local, cloneOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to get git repository: %s", err)
	}
	return repo, nil
}
