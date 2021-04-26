package file_explorer

import (
	"os"
	"path"
)

// ApplyToAllFiles applies a given function recursively to all files inside
// a given direrectory and its subdirectories
// returns an error
func ApplyToAllFiles(dir string, f func(string)) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			ApplyToAllFiles(file.Name(), f)
		} else {
			f(file.Name())
		}
	}
	return nil
}

// ApplyToAllFilesAsync applies a given function recursively and concurrently to all files inside
// a given direrectory and its subdirectories
// returns an error
func ApplyToAllFilesAsync(dir string, f func(string, chan bool), ch chan bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		ch <- true
		return err
	}
	childChan := make(chan bool, 2)
	childChan <- true
	for _, file := range files {
		if file.IsDir() {
			go ApplyToAllFilesAsync(path.Join(dir, file.Name()), f, childChan)
		} else {
			go f(path.Join(dir, file.Name()), childChan)
		}
		<-childChan
	}
	<-childChan
	ch <- true
	return nil
}

func ApplyToLeavesBeforeRoot(dir string, f func(string)) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			ApplyToLeavesBeforeRoot(path.Join(dir, file.Name()), f)
		}
		f(path.Join(dir, file.Name()))
	}
	return nil
}

func ApplyToLeavesBeforeRootAsync(dir string, f func(string, string, chan bool), quit chan bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		quit <- true
		return err
	}
	quitChild := make(chan bool, 2)
	quitChild <- true
	for _, file := range files {
		if file.IsDir() {
			ApplyToLeavesBeforeRootAsync(path.Join(dir, file.Name()), f, quitChild)
			<-quitChild
		}
		go f(dir, file.Name(), quitChild)
		<-quitChild
	}
	<-quitChild
	quit <- true
	return nil
}
