package jarfile

import (
	"archive/zip"
	"fmt"
	"g3tb-pngpacker/fileio"
	"os"
)

type JarFileReader struct {
	path string
}

func NewJarFileHandler(jarFilePath string) *JarFileReader {

	if _, err := os.Stat(jarFilePath); os.IsNotExist(err) {
		panic(jarFilePath + " not found")
	}

	j := &JarFileReader{
		path: jarFilePath,
	}
	return j
}

func (j *JarFileReader) ReadContentOf(fileInJar string) ([]byte, error) {
	jarFile, err := zip.OpenReader(j.path)
	if err != nil {
		return nil, err
	}
	defer jarFile.Close()

	fileInJarExists := false

	for _, f := range jarFile.File {
		if f.Name == fileInJar {
			fileInJarExists = true
			break
		}
	}

	if !fileInJarExists {
		return nil, fmt.Errorf("File '" + fileInJar + "' does not exist in " + j.path)
	}

	data, err := fileio.ReadBytes(fileInJar)
	if err != nil {
		return nil, err
	}
	return data, nil
}
