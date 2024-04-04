package svc

import (
	"io"
	"os"
)

type DiskFileStorageSvc struct {
}

func NewDiskFileStorageSvc() IFileStorageSvc {
	return &DiskFileStorageSvc{}
}

func (d *DiskFileStorageSvc) Read(ID string) (io.Reader, error) {
	f, err := os.Open(ID)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d *DiskFileStorageSvc) Write(ID string, r io.Reader) error {
	f, err := os.OpenFile(ID, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	b := make([]byte, 512)
	for {
		n, err := r.Read(b)
		f.Write(b[:n])
		if err == io.EOF {
			break
		}
	}
	return nil
}
