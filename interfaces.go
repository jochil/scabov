package vcs

//common interface for all vcs connectors
type Connector interface {
	Load(remote string, local string)
	AllCommits() map[string]*Commit
	AllDevelopers() map[string]*Developer
}

