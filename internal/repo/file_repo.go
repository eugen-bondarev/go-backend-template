package repo

type FileHandle int

type IFileRepo interface {
	GetPath(ID FileHandle) (string, error)
	AddPath(path string) (FileHandle, error)
}
