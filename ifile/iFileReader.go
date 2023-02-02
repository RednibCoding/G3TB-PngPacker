package ifile

import (
	"bytes"
	"fmt"
	"g3tb-pngpacker/fileio"
	"os"
	"path/filepath"
	"strconv"
)

// Png header: The first eight bytes of a PNG file are always the following values: 0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A
var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// The i file contains a charset (characters used in the game) followed by all the png files in the game
type IFileReader struct {
	// File name in the jar file
	fileName    string
	fileContent []byte
}

func NewIFileReader(data []byte) (*IFileReader, error) {
	// if !fileio.FileExist(path) {
	// 	return nil, fmt.Errorf(path + " not found")
	// }

	// data, err := fileio.ReadBytes(path)
	// if err != nil {
	// 	return nil, err
	// }

	iFile := IFileReader{fileName: "i", fileContent: data}
	return &iFile, nil
}

func (f *IFileReader) UnpackPngFilesFromJar(outputFolderPath string) (string, error) {
	offsets, err := f.findPngOffsetsInFile(f.fileContent)
	if err != nil {
		return "", err
	}
	charsetBuffer, pngBuffers, err := f.collectBuffersFromIFile(offsets, f.fileContent)
	if err != nil {
		return "", err
	}

	outputPath, err := fileio.MakeOrOverwriteDir(filepath.Join(outputFolderPath, "i_output"))
	if err != nil {
		return "", err
	}

	err = f.createFilesFromBuffers(charsetBuffer, pngBuffers, outputPath)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(pngBuffers)) + " png files created at " + outputPath, nil
}

func (f *IFileReader) findPngOffsetsInFile(data []byte) ([]int, error) {
	offsets := []int{}

	// copy the buffer as we don't want to change the original content of the file
	dataCpy := make([]byte, len(data))
	copy(dataCpy, data)

	// We need to remove bytes because "bytes.Index" always starts at the beginning of the buffer, so we would always find the first element
	numRemovedBytes := 0

	for {
		idx := bytes.Index(dataCpy, pngHeader)
		if idx == -1 {
			break
		}
		offsets = append(offsets, idx+numRemovedBytes)
		dataCpy = dataCpy[idx+len(pngHeader):]
		numRemovedBytes += idx + len(pngHeader)
	}
	return offsets, nil
}

func (f *IFileReader) collectBuffersFromIFile(offsets []int, data []byte) ([]byte, [][]byte, error) {
	var pngBuffers [][]byte
	last := -1

	if len(offsets) == 0 {
		return nil, nil, fmt.Errorf("this file does not contain any png data")
	}

	// append header (charset)
	charsetBuffer := data[0 : offsets[0]-1]

	if len(offsets) == 1 {
		pngBuffers = append(pngBuffers, data[offsets[0]:])
	} else {
		for _, idx := range offsets {
			if last == -1 {
				last = idx
				continue
			}
			pngBuffers = append(pngBuffers, data[last:idx-1])
			last = idx
		}
		if last < len(data) {
			pngBuffers = append(pngBuffers, data[last:])
		}
	}

	return charsetBuffer, pngBuffers, nil
}

func (f *IFileReader) createFilesFromBuffers(charsetBuffer []byte, pngBuffers [][]byte, outputFolderPath string) error {
	if len(pngBuffers) <= 0 || len(charsetBuffer) <= 0 {
		return fmt.Errorf(outputFolderPath + " not found")
	}

	if _, err := os.Stat(outputFolderPath); os.IsNotExist(err) {
		return err
	}

	// Write the (header) charset to a file (these vary from game version to game version)
	err := fileio.WriteBytes(filepath.Join(outputFolderPath, "charset.bin"), charsetBuffer)
	if err != nil {
		return err
	}

	fullPath := ""
	for i, buf := range pngBuffers {
		leadingZeros := len(strconv.Itoa(len(pngBuffers))) + 1
		format := "%0" + strconv.Itoa(leadingZeros) + "d"
		fullPath = filepath.Join(outputFolderPath, fmt.Sprintf(format, i)+".png")
		err := os.WriteFile(fullPath, buf, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
