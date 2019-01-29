package model

import (
	"errors"
	"io"
	"mime/multipart"
	"strings"

	tempfile "github.com/mash/go-tempfile-suffix"
)

type File struct {
	File       multipart.File
	Header     *multipart.FileHeader
	Prefix     string
	ValidTypes map[string]string
	Path       string
}

func (f *File) SaveFile() (string, error) {

	file, err := tempfile.TempFileWithSuffix(f.Path, f.Prefix, "."+f.getExtensionFromName())
	if err != nil {
		return "", err
	}
	defer file.Close()
	io.Copy(file, f.File)

	dirs := strings.Split(file.Name(), "/")

	return dirs[len(dirs)-1], nil
}

func (f *File) getExtensionFromName() string {
	name := strings.Split(f.Header.Filename, ".")
	return name[1]
}

func (f *File) Validate() error {

	for extension, mime := range f.ValidTypes {
		if mime == f.Header.Header["Content-Type"][0] && extension == f.getExtensionFromName() {
			return nil
		}
	}

	return errors.New("invalid file type provided")
}
