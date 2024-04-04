package svc

import "io"

type IFileStorageSvc interface {
	Read(ID string) (io.Reader, error)
	Write(ID string, r io.Reader) error
}
