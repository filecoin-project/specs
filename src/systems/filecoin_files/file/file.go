package file

import (
	"io/ioutil"
)

func ReadAll(f File) ([]byte, error) {
	f2 := FileReadWriter{f}
	return ioutil.ReadAll(f2)
}

type FileReadWriter struct {
	f File
}

func (f FileReadWriter) Read(buf []byte) (n int, err error) {
	ret := f.f.Read(buf)
	return ret.size(), ret.e()
}

func (f FileReadWriter) Write(buf []byte) (n int, err error) {
	ret := f.f.Write(buf)
	return ret.size(), ret.e()
}
