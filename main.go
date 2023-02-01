package main

import (
	"fmt"
	"g3tb-pngpacker/fileio"
	"g3tb-pngpacker/ifile"
	"os"
	"path/filepath"
)

/*
G3TB-PngPacker by Michael Binder

Tool for extracting and packing png images from the game Gothic 3 The Beginning
*/

func main() {
	printVersion()

	dragNdropFile, exeDir, err := getArgs()
	if err != nil {
		panic(err)
	}

	// dragNdropFile := "C:\\Users\\mlb\\Desktop\\pnpteset\\i_output"
	// exeDir := "C:\\Users\\mlb\\Desktop\\pnpteset\\"

	if !fileio.FileExist(dragNdropFile) {
		panic("File not found: " + dragNdropFile)
	}

	if fileio.IsDirectory(dragNdropFile) {
		iFileWriter, err := ifile.NewIFileWriter(dragNdropFile)
		if err != nil {
			panic(err)
		}
		result, err := iFileWriter.PackPngFilesIntoIFile(exeDir)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	} else {
		// if filepath.Ext(dragNdropFile) != ".jar" {
		// 	panic("invalid file, please provide the Gothic 3 The Beginnig jar file or the folder with png files")
		// }
		iFileReader, err := ifile.NewIFileReader(dragNdropFile)
		if err != nil {
			panic(err)
		}

		result, err := iFileReader.UnpackPngFilesFromIFile(exeDir)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}

}

func getArgs() (string, string, error) {
	var dragNdropFile = ""

	if len(os.Args) < 2 {
		return "", "", fmt.Errorf("no file argument provided, please provide the Gothic 3 The Beginnig jar file or the folder with png files")
	}

	dragNdropFile = os.Args[1]
	exeDir := filepath.Dir(os.Args[0])

	return dragNdropFile, exeDir, nil
}

func printVersion() {
	version := "1.0.1"
	fmt.Println("")
	fmt.Println("###########################################################")
	fmt.Println("#     G3TB-PngPacker version " + version + " by Michael Binder      #")
	fmt.Println("###########################################################")
	fmt.Println("")
	fmt.Println("Found any bugs? Please create a new issue at: https://github.com/RednibCoding/G3TB-PngPacker/issues")
	fmt.Println("")
}
