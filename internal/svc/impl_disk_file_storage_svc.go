package svc

import (
	"fmt"
	"io"
	"os"
)

type DiskFileStorageSvc struct {
	location string
}

func NewDiskFileStorage(location string) IFileStorageSvc {
	return &DiskFileStorageSvc{
		location: location,
	}
}

func (d *DiskFileStorageSvc) getPath(name string) string {
	return fmt.Sprintf("%s/%s", d.location, name)
}

func (d *DiskFileStorageSvc) Read(name string) (io.Reader, error) {
	f, err := os.Open(d.getPath(name))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d *DiskFileStorageSvc) Write(name string, r io.Reader) error {
	f, err := os.OpenFile(d.getPath(name), os.O_RDWR|os.O_CREATE, 0644)
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
