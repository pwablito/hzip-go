package input

import "io/fs"

type Meta interface {
	GetMode() fs.FileMode
}

type FileMeta struct {
	Mode     fs.FileMode
	Owner_ID int
	Group_ID int
}

func (meta FileMeta) GetMode() fs.FileMode {
	return meta.Mode
}

type DirectoryMeta struct {
	Mode     fs.FileMode
	Owner_ID int
	Group_ID int
	Name     string
}

func (meta DirectoryMeta) GetMode() fs.FileMode {
	return meta.Mode
}
