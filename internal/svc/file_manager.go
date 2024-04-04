package svc

import (
	"go-backend-template/internal/repo"
	"io"
)

type IFileManager interface {
	Read(ID repo.FileHandle) (io.Reader, error)
	Write(name string, r io.Reader) (repo.FileHandle, error)
}
