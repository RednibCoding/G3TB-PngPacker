package main

import (
	"fmt"
	"g3tb-pngpacker/fileio"
	"g3tb-pngpacker/ifile"
	"g3tb-pngpacker/jarfile"
	"os"
	"path/filepath"
)

/*
G3TB-PngPacker by Michael Binder

Tool for extracting and packing png images from the game Gothic 3 The Beginning
*/

func main() {
	printVersion()

	dragNdropFile, dragNdropJar, exeDir, err := parseArgs()
	if err != nil {
		panicwait(err)
	}

	// dragNdropFile := "C:\\Users\\mlb\\Desktop\\pnpteset\\i_output"
	// dragNdropFile = ""
	// dragNdropJar := "C:\\Users\\mlb\\Desktop\\pnpteset\\tb - Kopie.jar"

	// exeDir := "C:\\Users\\mlb\\Documents\\DEV\\github\\go\\G3TB-PngPacker\\bin"

	// if !fileio.FileExist(dragNdropFile) {
	// 	panic("File not found: " + dragNdropFile)
	// }

	// if dragNdropJar == "" {
	// 	panic(invalidArgsStr())
	// }

	println("folder: " + dragNdropFile)
	println("jar: " + dragNdropJar)
	println("exe: " + exeDir)
	println("")
	// waitkey()

	if !fileio.IsDirectory(dragNdropFile) {

		jarFileReader, err := jarfile.NewJarFileReader(dragNdropJar)
		if err != nil {
			panicwait(err)
		}

		data, err := jarFileReader.ReadContentOf("i")
		if err != nil {
			panicwait(err)
		}

		iFileReader, err := ifile.NewIFileReader(data)
		if err != nil {
			panicwait(err)
		}

		result, err := iFileReader.UnpackPngFilesFromJar(filepath.Dir(dragNdropJar))
		if err != nil {
			panicwait(err)
		}
		printwait(result)
	} else {

		iFileWriter, err := ifile.NewIFileWriter(dragNdropFile)
		if err != nil {
			panicwait(err)
		}

		result, err := iFileWriter.PackPngFilesIntoJar(dragNdropJar)
		if err != nil {
			panicwait(err)
		}
		printwait(result)
	}
}

func parseArgs() (string, string, string, error) {
	if len(os.Args) == 2 {
		// it must the jar file
		if filepath.Ext(os.Args[1]) != ".jar" {
			return "", "", "", fmt.Errorf(invalidArgsStr())
		}
		return "", os.Args[1], filepath.Dir(os.Args[0]), nil
	}

	if len(os.Args) == 3 {
		// it must be a folder and the jar file
		if filepath.Ext(os.Args[1]) == ".jar" && fileio.IsDirectory(os.Args[2]) {
			return os.Args[2], os.Args[1], filepath.Dir(os.Args[0]), nil
		} else if filepath.Ext(os.Args[2]) == ".jar" && fileio.IsDirectory(os.Args[1]) {
			return os.Args[1], os.Args[2], filepath.Dir(os.Args[0]), nil
		}
	}

	return "", "", "", fmt.Errorf(invalidArgsStr())
}

func printVersion() {
	version := "1.1.0"
	fmt.Println("")
	fmt.Println("###########################################################")
	fmt.Println("#     G3TB-PngPacker version " + version + " by Michael Binder      #")
	fmt.Println("###########################################################")
	fmt.Println("")
	fmt.Println("Found any bugs? Please create a new issue at: https://github.com/RednibCoding/G3TB-PngPacker/issues")
	fmt.Println("")
}

func invalidArgsStr() string {
	return ("invalid file: either provide the Gothic 3 The Beginnig jar alone for extracting the png files, or the Gothic 3 The Beginnig jar file AND the folder with png files to pack the folder into the provided jar file")
}

func panicwait(err error) {
	fmt.Println(err)
	waitkey()
}

func printwait(msg string) {
	fmt.Println(msg)
	waitkey()
}
func waitkey() {
	os.Stdin.Read(make([]byte, 1))
}
