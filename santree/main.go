package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"timmy.narnian.us/git/timmy/sanitize"
)

func main() {
	var previousDir string = "./"
	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		fmt.Printf("%s >>>>> %s\t%t\n", previousDir, path, strings.HasPrefix(path, previousDir) || previousDir == ".")
		defer func() { previousDir = filepath.Dir(path) }()
		if err != nil {
			println(err)
		}

		if info.IsDir() && !DirIsEmpty(path) {
			return nil
		}

		newpath := sanitize.SanitizeFilepath(filepath.Dir(path))
		filename := sanitize.SanitizeFilename(info.Name())
		if path != filepath.Join(newpath, filename) {
			fmt.Println(path, "->", filepath.Join(newpath, filename))
			os.MkdirAll(newpath, os.ModePerm)
			fmt.Println(os.Rename(path, filepath.Join(newpath, filename)))
		}
		if !(strings.HasPrefix(path, previousDir) || previousDir == "./") {
			rmpath := previousDir
			for len(rmpath) > 1 {
				fmt.Println("rmdir", rmpath)
				fmt.Println(rmpath, os.Remove(rmpath))
				rmpath = filepath.Dir(rmpath)
			}
		}

		return nil
	})
	fmt.Println(previousDir)
}

// also returns false if os.Open fails
func DirIsEmpty(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	if _, err = file.Readdirnames(1); err == io.EOF {
		return true
	}
	return false
}
