package fileDetection

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsSymlink(file string) bool {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		return false
	}
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	} else {
		return false
	}
}

func ListDir(searchDir string) []string {
	var fileList []string

	// edode : Mapping all the files from the root directory of the OS
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !IsDir(path) {
			fileList = append(fileList, path)
		}
		return err
	})
	if e != nil {
		panic(e)
	}

	return fileList
}

func ListRootDir(root string) []string {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files
}

func GetAllFiles(filepath string) []string {
	if filepath == "" {
		filepath = "/"
	}
	var allFiles []string
	for _, dirs := range ListRootDir(filepath) {
		if len(allFiles) == 0 {
			allFiles = ListDir(filepath + dirs)
		} else {
			temp := make([]string, len(allFiles))
			copy(temp, allFiles)
			allFiles = ListDir(filepath + dirs)
			for _, file := range temp {
				allFiles = append(allFiles, file)
			}
		}
	}
	return allFiles
}

func IsDir(path string) bool {
	name := path
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}
