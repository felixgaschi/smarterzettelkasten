package file_explorer

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

// ApplyToAllFilesAsync applies a given function recursively and concurrently to all files inside
// a given direrectory and its subdirectories
// returns an error
func ApplyToAllFilesAsync(dir string, f func(string, string, chan bool), ch chan bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		ch <- true
		return err
	}

	childChan := make(chan bool, len(files))
	for _, file := range files {
		if file.IsDir() {
			go ApplyToAllFilesAsync(path.Join(dir, file.Name()), f, childChan)
		} else {
			go f(dir, file.Name(), childChan)
		}
	}
	for range files {
		<-childChan
	}
	ch <- true
	return nil
}

// ApplyToAllFilesAsync applies a given function recursively and concurrently to all files inside
// a given direrectory and its subdirectories and it insures to run the script on children directories
// and files before parents
func ApplyToLeavesBeforeRootAsync(dir string, f func(string, string, chan bool), quit chan bool) {
	defer func(quit chan bool) {
		quit <- true
	}(quit)

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	quitChild := make(chan bool, len(files))
	for _, file := range files {
		go func(dir string, f func(string, string, chan bool), second_quit chan bool, file fs.DirEntry) {
			if file.IsDir() {
				first_quit := make(chan bool)
				go ApplyToLeavesBeforeRootAsync(path.Join(dir, file.Name()), f, first_quit)
				<-first_quit
			}
			go f(dir, file.Name(), second_quit)
		}(dir, f, quitChild, file)
	}
	for range files {
		<-quitChild
	}
}
