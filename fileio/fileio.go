package fileio

import "os"

func ReadBytes(filePath string) ([]byte, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func WriteBytes(filePath string, data []byte) error {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetAllFilesInDir(directoryPath string) ([]string, error) {
	dir, err := os.Open(directoryPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	fileNames := []string{}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeOrOverwriteDir(path string) (string, error) {
	if FileExist(path) {
		err := os.RemoveAll(path)
		if err != nil {
			return "", err
		}
	}
	err := os.Mkdir(path, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}
