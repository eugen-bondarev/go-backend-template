package svc

import "io"

type IFileStorage interface {
	Read(name string) (io.Reader, error)
	Write(name string, r io.Reader) error
}
