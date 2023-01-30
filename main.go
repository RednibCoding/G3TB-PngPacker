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
PngPacker by Michael Binder

PngPacker is a tool that scans the given file for png signatures and if given, extracts those png files.
It also can pack a given folder of png files into a pack file. While packing, it writes a PngPacker signature
so it knows when unpacking, that it is a file packed by PngPacker. The signature is an arbitrary sequence of five bytes at the beginning see: getPngPackerSignature()
Also it packs the png file names in front of the png signatures to preserve the original file names of the png files.
So the file structure of a packed file looks like as follows:
 --------------------
| PngPackerSignature |  -> the signature to identify wether it is a packed file from PngPacker
| #myImage1.png#     |  -> original file name of the first image
| png signature      |  -> png signature of the first image
| png content        |  -> png content of the first image
| #myImage2.png#     |  -> original file name of the second image
| png signature      |  -> png signature of the second image
| png content        |  -> png content of the second image
| ...                |
 --------------------

 PngPacker can not only unpack png images from its own produced pack files but any file that contains png signatures and content.
 PngPacker can identify automatically wether it is a file packed by itself or a random other file.
*/

var output_postfix = "_output"
var g3tb_charset = []byte{
	// Offset 0x00000000 to 0x00000145
	0x13, 0xB4, 0x00, 0x9E, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03,
	0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x64, 0x00, 0x65,
	0x00, 0x66, 0x00, 0x68, 0x00, 0x6A, 0x00, 0x6D, 0x00, 0x71, 0x00, 0x72,
	0x00, 0x74, 0x00, 0x78, 0x00, 0x79, 0x00, 0x7B, 0x00, 0x7E, 0x00, 0xC8,
	0x00, 0xC9, 0x00, 0xCA, 0x00, 0xCB, 0x00, 0xCC, 0x00, 0xCD, 0x00, 0xCF,
	0x00, 0xD3, 0x00, 0xD4, 0x00, 0xD5, 0x00, 0xDD, 0x00, 0xDE, 0x00, 0xE0,
	0x00, 0xE7, 0x00, 0xE8, 0x00, 0xE9, 0x00, 0xEA, 0x00, 0xEB, 0x00, 0xEC,
	0x00, 0xF1, 0x00, 0xF2, 0x00, 0xF3, 0x00, 0xF4, 0x00, 0xFA, 0x00, 0xFB,
	0x00, 0xFC, 0x00, 0xFD, 0x01, 0x04, 0x01, 0x05, 0x01, 0x0E, 0x01, 0x0F,
	0x01, 0x10, 0x01, 0x2C, 0x01, 0x2D, 0x01, 0x2E, 0x01, 0x2F, 0x01, 0x30,
	0x01, 0x32, 0x01, 0x33, 0x01, 0x34, 0x01, 0x35, 0x01, 0x37, 0x01, 0x38,
	0x01, 0x3B, 0x01, 0x3C, 0x01, 0x3D, 0x01, 0x3E, 0x01, 0x3F, 0x01, 0x40,
	0x01, 0x43, 0x01, 0x44, 0x01, 0x45, 0x01, 0x46, 0x01, 0x49, 0x01, 0x4B,
	0x01, 0x4C, 0x01, 0x4D, 0x01, 0x4E, 0x01, 0x4F, 0x01, 0x50, 0x01, 0x52,
	0x01, 0x53, 0x01, 0x54, 0x01, 0x55, 0x01, 0x56, 0x01, 0x57, 0x01, 0x58,
	0x01, 0x5A, 0x01, 0x5C, 0x01, 0x5F, 0x01, 0x60, 0x01, 0x61, 0x01, 0x62,
	0x01, 0x64, 0x01, 0x66, 0x01, 0x67, 0x01, 0x68, 0x01, 0x69, 0x01, 0x6B,
	0x01, 0x6C, 0x01, 0x6D, 0x01, 0x6E, 0x01, 0x71, 0x01, 0x76, 0x01, 0x78,
	0x01, 0x79, 0x01, 0x7A, 0x01, 0x7B, 0x01, 0x92, 0x01, 0x98, 0x01, 0xF5,
	0x02, 0x2C, 0x02, 0x62, 0x02, 0x6D, 0x02, 0x76, 0x02, 0x9B, 0x02, 0x9C,
	0x02, 0xC0, 0x02, 0xC1, 0x02, 0xC5, 0x02, 0xC7, 0x02, 0xC8, 0x02, 0xC9,
	0x02, 0xCA, 0x02, 0xCB, 0x03, 0x20, 0x03, 0x21, 0x03, 0x22, 0x03, 0x23,
	0x03, 0x24, 0x03, 0x26, 0x03, 0x27, 0x03, 0x28, 0x03, 0x29, 0x03, 0x2A,
	0x03, 0x2B, 0x03, 0x2C, 0x03, 0x2D, 0x03, 0x2E, 0x03, 0x30, 0x03, 0x31,
	0x03, 0x33, 0x03, 0x34, 0x03, 0x35, 0x03, 0x36, 0x03, 0x37, 0x03, 0x38,
	0x03, 0x3A, 0x03, 0x3B, 0x03, 0x52, 0x03, 0x53, 0x03, 0x54, 0x03, 0x55,
	0x03, 0x56, 0x03, 0x57, 0x03, 0x58, 0x03, 0x59, 0x00, 0x00, 0x13, 0xB4,
	0x00, 0x00,
}

func main() {
	var dndFile = ""
	if len(os.Args) >= 2 {
		dndFile = os.Args[1]
	}

	dndFile = "C:\\Users\\mlb\\Documents\\DEV\\github\\go\\G3TB-PngPacker\\j_1_output"

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

	mergedBuffer = append(mergedBuffer, g3tb_charset...)

	for _, buffer := range buffers {
		mergedBuffer = append(mergedBuffer, buffer...)
	}

	fullPath := filepath.Join(path, outputFileName)
	fullPath = strings.Replace(fullPath, output_postfix, "", -1)
	err := os.WriteFile(fullPath, mergedBuffer, 0644)
	if err != nil {
		waitExit(err.Error())
	}

	waitExit(strconv.Itoa(len(buffers)) + " png files packed into " + fullPath)
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

	fmt.Println(strconv.Itoa(len(offsets)) + " PNG images found")
	return offsets
}

func collectPngBuffers(offsets []int, data []byte) [][]byte {
	var buffers [][]byte
	last := -1

	if len(offsets) == 0 {
		waitExit("This file does not contain any png data")
	}

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

	for i, buf := range buffers {
		fileName := "image_"
		fullPath = filepath.Join(path+output_postfix, fileName+strconv.Itoa(i)+".png")
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
