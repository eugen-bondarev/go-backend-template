package svc

import "io"

type IFileStorageSvc interface {
	Read(name string) (io.Reader, error)
	Write(name string, r io.Reader) error
}
