package upload

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
)

type File struct {
	Name    string
	Content []byte
}

func Unzip(zipData []byte) ([]File, error) {
	reader := bytes.NewReader(zipData)
	zipReader, err := zip.NewReader(reader, int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	var files []File

	for _, f := range zipReader.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		content, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		files = append(files, File{
			Name:    f.Name,
			Content: content,
		})
	}

	return files, nil
}
