package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
G3TB-PngPacker by Michael Binder

Tool for extracting and packing png images from the game Gothic 3 The Beginning
*/

var output_postfix = "_output"

func main() {
	var dndFile = ""
	if len(os.Args) >= 2 {
		dndFile = os.Args[1]
	}

	// dndFile = "C:\\Users\\mlb\\Documents\\DEV\\github\\go\\G3TB-PngPacker\\j_my_output"

	if dndFile == "" {
		waitExit("Please drag&drop the file onto the 'PngPacker' executable")
	}

	fmt.Println("Processing: " + dndFile + " ...")

	if _, err := os.Stat(dndFile); os.IsNotExist(err) {
		waitExit(dndFile + " not found")
	}

	if isDirectory(dndFile) {
		packPngs(dndFile)
	} else {
		unpackPngs(dndFile)
	}
}

// --- Packing ---

func packPngs(path string) {
	pngPaths := collectFileNamesInDir(path)
	pngBuffers := createPngBuffersFromPngFiles(pngPaths)

	// Go up one directory as we want to create the packfile at the same location where the user has his folder with png files
	dir := filepath.Base(path)
	updir := path[0 : len(path)-len(dir)-1] // +1 to also remove the slash
	writePngBuffersAsPackFile(updir, dir, pngBuffers)
}

func collectFileNamesInDir(path string) []string {
	pngNames := make([]string, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		waitExit(err.Error())
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {
			pngNames = append(pngNames, filepath.Join(path, file.Name()))
		}
	}
	return pngNames
}

func createPngBuffersFromPngFiles(pngFilePaths []string) [][]byte {
	if len(pngFilePaths) == 0 {
		waitExit("No png files found in directory")
	}

	buffers := make([][]byte, 0)

	// Read header (charset) file
	charsetFilePath := filepath.Join(filepath.Dir(pngFilePaths[0]), "charset")
	if _, err := os.Stat(charsetFilePath); os.IsNotExist(err) {
		waitExit(charsetFilePath + " not found")
	}
	charsetBuf := readBytes(charsetFilePath)
	buffers = append(buffers, charsetBuf)

	for _, path := range pngFilePaths {
		buffer := readBytes(path)
		// skip empty files
		if len(buffer) == 0 {
			continue
		}
		buffers = append(buffers, buffer)
	}
	return buffers
}

func writePngBuffersAsPackFile(path string, outputFileName string, buffers [][]byte) {
	if !isDirectory(path) {
		waitExit(path + " is not a valid output path")
	}

	mergedBuffer := make([]byte, 0)

	// mergedBuffer = append(mergedBuffer, g3tb_charset...)

	for i, buffer := range buffers {
		mergedBuffer = append(mergedBuffer, buffer...)
		if i < len(buffers)-1 {
			mergedBuffer = append(mergedBuffer, byte(0x00))
		}
	}

	fullPath := filepath.Join(path, outputFileName)
	fullPath = strings.Replace(fullPath, output_postfix, "", -1)
	err := os.WriteFile(fullPath, mergedBuffer, 0644)
	if err != nil {
		waitExit(err.Error())
	}

	waitExit(strconv.Itoa(len(buffers)-1) + " png files packed into " + fullPath)
}

// --- Unpacking ---

func unpackPngs(filePath string) {
	data := readBytes(filePath)
	var fileNames []string
	offsets := findPngOffsets(data)

	pngBuffers := collectPngBuffers(offsets, data)
	writePngBuffers(pngBuffers, fileNames, filepath.Dir(filePath+"\\"))
}

func findPngOffsets(data []byte) []int {
	pngStartPattern := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	offsets := []int{}
	dataCpy := make([]byte, len(data))
	numRemoved := 0
	copy(dataCpy, data)

	for {
		idx := bytes.Index(dataCpy, pngStartPattern)
		if idx == -1 {
			break
		}
		offsets = append(offsets, idx+numRemoved)
		dataCpy = dataCpy[idx+len(pngStartPattern):]
		numRemoved += idx + len(pngStartPattern)
	}

	fmt.Println(strconv.Itoa(len(offsets)) + " png image\\s found")
	return offsets
}

func collectPngBuffers(offsets []int, data []byte) [][]byte {
	var buffers [][]byte
	last := -1

	if len(offsets) == 0 {
		waitExit("This file does not contain any png data")
	}

	// append header (charset)
	buffers = append(buffers, data[0:offsets[0]-1])

	if len(offsets) == 1 {
		buffers = append(buffers, data[offsets[0]:])
	} else {
		for _, idx := range offsets {
			if last == -1 {
				last = idx
				continue
			}
			buffers = append(buffers, data[last:idx-1])
			last = idx
		}
		if last < len(data) {
			buffers = append(buffers, data[last:])
		}
	}

	return buffers
}

func writePngBuffers(buffers [][]byte, fileNames []string, path string) {
	fullPath := ""
	if len(buffers) > 0 {
		if _, err := os.Stat(path + output_postfix); os.IsNotExist(err) {
			err := os.Mkdir(path+output_postfix, os.ModePerm)
			if err != nil {
				waitExit(err.Error())
			}
		}
	}

	// Write the (header) charset to a file (these vary from game version to game version)
	charset := buffers[0]
	writeBytes(filepath.Join(path+output_postfix, "charset"), charset)
	buffers = buffers[1:]

	for i, buf := range buffers {
		fileName := "image_"
		leadingZeros := len(strconv.Itoa(len(buffers))) + 1
		format := "%0" + strconv.Itoa(leadingZeros) + "d"
		fullPath = filepath.Join(path+output_postfix, fileName+fmt.Sprintf(format, i)+".png")
		err := os.WriteFile(fullPath, buf, 0644)
		if err != nil {
			waitExit(err.Error())
		}
	}
	waitExit(strconv.Itoa(len(buffers)) + " png files written to: " + filepath.Dir(fullPath))
}

// --- Helpers ---

func readBytes(filePath string) []byte {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		waitExit(err.Error())
	}
	return buf
}

func writeBytes(filePath string, data []byte) {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		waitExit(err.Error())
	}
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		waitExit(err.Error())
	}

	return fileInfo.IsDir()
}

func waitExit(msg ...string) {
	if len(msg) > 0 {
		if msg[0] != "" {
			fmt.Println(msg)
		}
	}
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	os.Exit(0)
}

// func reverseSlice[T comparable](s []T) {
// 	sort.SliceStable(s, func(i, j int) bool {
// 		return i > j
// 	})
// }

// func dumpBuffers(buffers [][]byte) {
// 	for _, v := range buffers {
// 		for _, d := range v {
// 			fmt.Printf("0x%02X ", d)
// 		}
// 	}
// }
