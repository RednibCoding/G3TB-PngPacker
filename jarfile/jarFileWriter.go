package jarfile

import (
	"archive/zip"
	"bytes"
	"fmt"
	"g3tb-pngpacker/fileio"
	"io"
	"os"
	"path/filepath"
)

type JarFileWriter struct {
	jarFilePath string
}

func NewJarFileWriter(jarFilePath string) (*JarFileWriter, error) {

	if !fileio.FileExist(jarFilePath) {
		return nil, fmt.Errorf(jarFilePath + " not found")
	}

	j := &JarFileWriter{
		jarFilePath: jarFilePath,
	}
	return j, nil
}

func (j *JarFileWriter) WriteFileToJar(fileNameToWriteNew string, newFileContent []byte) error {
	// Open the original jar file
	r, err := zip.OpenReader(j.jarFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create a new buffer to hold the contents of the updated jar file
	var buf bytes.Buffer

	// Create a new zip writer to write to the buffer
	w := zip.NewWriter(&buf)

	// Replace the contents of the file named 'fileNameToWriteNew'
	for _, f := range r.File {
		if f.Name == fileNameToWriteNew {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			fw, err := w.Create(f.Name)
			if err != nil {
				return err
			}

			// Write the new contents to the zip writer
			_, err = io.Copy(fw, bytes.NewBuffer(newFileContent))
			if err != nil {
				return err
			}

			continue
		}

		frc, err := f.Open()
		if err != nil {
			return err
		}
		defer frc.Close()

		fw, err := w.Create(f.Name)
		if err != nil {
			return err
		}

		if _, err = io.Copy(fw, frc); err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	// Create a new file for the updated jar
	f, err := os.Create(filepath.Join(filepath.Dir(j.jarFilePath), "tmp.jar"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the contents of the buffer to the new file
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}

	f.Close()
	r.Close()

	// Replace the old JAR with the new JAR
	err = os.Rename(filepath.Join(filepath.Dir(j.jarFilePath), "tmp.jar"), j.jarFilePath)
	if err != nil {
		return err
	}

	return nil
}
