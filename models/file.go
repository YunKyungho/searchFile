package models

type File struct {
	Name    string
	ModTime string
}

type Directory struct {
	ModTime string
	Child   []File
	DiNo    int
}
