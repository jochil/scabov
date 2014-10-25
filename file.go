package vcs

type File struct {
	Id       string
	Path     string
	Size     int64
	Contents []byte
}

func (f *File) String() string {
	return string(f.Contents[:])
}
