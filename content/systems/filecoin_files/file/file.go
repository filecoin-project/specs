package file

import (
	"io/ioutil"
)

func ReadAll(f File) ([]byte, error) {
	f2 := FileReadWriter{f, 0}
	return ioutil.ReadAll(f2)
}

type FileReadWriter struct {
	f      File
	offset int
}

func (f FileReadWriter) Read(buf []byte) (n int, err error) {
	ret := f.f.Read(f.offset, len(buf), buf)
	f.offset += ret.size()
	return ret.size(), ret.e()
}

func (f FileReadWriter) Write(buf []byte) (n int, err error) {
	ret := f.f.Write(f.offset, len(buf), buf)
	f.offset += ret.size()
	return ret.size(), ret.e()
}

func FromPath(path Path) *FileReadWriter {
	return &FileReadWriter{} // TODO: move to using Filestore
}
