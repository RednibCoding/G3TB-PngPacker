package jarfile

import (
	"archive/zip"
	"bytes"
	"fmt"
	"g3tb-pngpacker/fileio"
	"io"
)

type JarFileReader struct {
	jarFilePath string
}

func NewJarFileReader(jarFilePath string) (*JarFileReader, error) {

	if !fileio.FileExist(jarFilePath) {
		return nil, fmt.Errorf(jarFilePath + " not found")
	}

	j := &JarFileReader{
		jarFilePath: jarFilePath,
	}
	return j, nil
}

func (j *JarFileReader) ReadContentOf(fileInJar string) ([]byte, error) {
	jarFile, err := zip.OpenReader(j.jarFilePath)
	if err != nil {
		return nil, err
	}
	defer jarFile.Close()

	var iFileReader io.ReadCloser

	for _, f := range jarFile.File {
		if f.Name == fileInJar {
			iFileReader, err = f.Open()
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if iFileReader == nil {
		return nil, fmt.Errorf("File '" + fileInJar + "' does not exist in " + j.jarFilePath)
	}

	var data bytes.Buffer

	_, err = io.Copy(&data, iFileReader)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}
