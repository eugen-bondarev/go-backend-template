package svc

import (
	"go-backend-template/internal/repo"
	"io"
)

type FileManager struct {
	fileRepo    repo.IFileRepo
	fileStorage IFileStorageSvc
}

func NewFileManager(fileRepo repo.IFileRepo, fileStorage IFileStorageSvc) IFileManager {
	return &FileManager{
		fileRepo:    fileRepo,
		fileStorage: fileStorage,
	}
}

func (fm *FileManager) Read(ID repo.FileHandle) (io.Reader, error) {
	path, err := fm.fileRepo.GetPath(ID)
	if err != nil {
		return nil, err
	}

	return fm.fileStorage.Read(path)
}

func (fm *FileManager) Write(name string, r io.Reader) (repo.FileHandle, error) {
	handle, err := fm.fileRepo.AddPath(name)
	if err != nil {
		return handle, err
	}

	return handle, fm.fileStorage.Write(name, r)
}
