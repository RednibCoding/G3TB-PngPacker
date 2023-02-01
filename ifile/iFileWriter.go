package ifile

import (
	"fmt"
	"g3tb-pngpacker/fileio"
	"os"
	"path/filepath"
	"strconv"
)

type IFileWriter struct {
	// File name in the jar file
	FileName      string
	pngFolderPath string
}

func NewIFileWriter(pngFolderPath string) (*IFileWriter, error) {
	if !fileio.FileExist(pngFolderPath) {
		return nil, fmt.Errorf(pngFolderPath + " not found")
	}
	iFile := IFileWriter{FileName: "i", pngFolderPath: pngFolderPath}
	return &iFile, nil
}

func (f *IFileWriter) PackPngFilesIntoIFile(outputPath string) (string, error) {
	pngPaths, err := f.getPngFiles()
	if err != nil {
		return "", err
	}

	charsetPath, err := f.getCharsetFile()
	if err != nil {
		return "", err
	}

	pngBuffers, err := f.createPngBuffersFromPngFiles(pngPaths)
	if err != nil {
		return "", err
	}

	charsetBuffer, err := f.createCharsetBufferFromCharsetFile(charsetPath)
	if err != nil {
		return "", err
	}

	err = f.createIFileFromBuffers(outputPath, charsetBuffer, pngBuffers)
	if err != nil {
		return "", err
	}

	return "i file with " + strconv.Itoa(len(pngBuffers)) + " png files created at " + outputPath, nil
}

func (f *IFileWriter) getPngFiles() ([]string, error) {
	files, err := fileio.GetAllFilesInDir(f.pngFolderPath)
	if err != nil {
		return nil, err
	}

	pngFiles := []string{}
	for _, file := range files {
		if filepath.Ext(file) == ".png" {
			pngFiles = append(pngFiles, filepath.Join(f.pngFolderPath, file))
		}
	}
	return pngFiles, nil
}

func (f *IFileWriter) getCharsetFile() (string, error) {
	files, err := fileio.GetAllFilesInDir(f.pngFolderPath)
	if err != nil {
		return "", err
	}

	charsetFile := ""
	for _, file := range files {
		if filepath.Base(file) == "charset.bin" {
			charsetFile = filepath.Join(f.pngFolderPath, file)
			break
		}
	}
	return charsetFile, nil
}

func (f *IFileWriter) createPngBuffersFromPngFiles(pngFiles []string) ([][]byte, error) {
	if len(pngFiles) == 0 {
		return nil, fmt.Errorf("no png files found in directory")
	}

	buffers := make([][]byte, 0)

	for _, path := range pngFiles {
		buffer, err := fileio.ReadBytes(path)
		if err != nil {
			return nil, err
		}
		// skip empty files
		if len(buffer) == 0 {
			continue
		}
		buffers = append(buffers, buffer)
	}
	return buffers, nil
}

func (f *IFileWriter) createCharsetBufferFromCharsetFile(charsetFile string) ([]byte, error) {
	charsetBuf, err := fileio.ReadBytes(charsetFile)
	if err != nil {
		return nil, err
	}
	charsetBuf = append(charsetBuf, 0x00)
	return charsetBuf, nil
}

func (f *IFileWriter) createIFileFromBuffers(outputFilePath string, charsetBuffer []byte, pngBuffers [][]byte) error {
	mergedBuffer := make([]byte, 0)
	mergedBuffer = append(mergedBuffer, charsetBuffer...)

	for i, buffer := range pngBuffers {
		mergedBuffer = append(mergedBuffer, buffer...)
		if i < len(pngBuffers)-1 {
			mergedBuffer = append(mergedBuffer, byte(0x00))
		}
	}

	outputPath := filepath.Join(outputFilePath, f.FileName)

	err := os.WriteFile(outputPath, mergedBuffer, 0644)

	if err != nil {
		return err
	}

	return nil
}
