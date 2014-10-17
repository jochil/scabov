package vcs

//common interface for all vcs connectors
type Connector interface {
	Load(remote string, local string)
	Commits() map[string]*Commit
	Developers() map[string]*Developer
}

