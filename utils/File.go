package utils

import "os"

type ReceivedFile struct {
	Pointer *os.File
	Name    string
	Path    string
}

func (file *ReceivedFile) Close() {
	file.Pointer.Close()
	file = nil
}
