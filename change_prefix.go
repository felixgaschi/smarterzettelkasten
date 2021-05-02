package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/felixgaschi/smarterzettelkasten/fileexplorer"
	log "github.com/sirupsen/logrus"
)

// Return a closure that will replace a prefix for a given file
// send the new and old title in a channel to follow change
func ChangeFileAndDirPrefix(oldPrefix, newPrefix string, newTitles chan [2]string) func(string, string, chan bool) {
	return func(base, name string, quit chan bool) {
		defer SendTrue(quit)
		if strings.HasPrefix(name, oldPrefix) {
			newFname := newPrefix + strings.TrimPrefix(name, oldPrefix)
			oldpath := path.Join(base, name)
			newpath := path.Join(base, newFname)
			err := os.Rename(oldpath, newpath)
			if err != nil {
				fmt.Println(err)
				log.Error(fmt.Sprintf("Could not rename '%s' into '%s'", oldpath, newpath))
				return
			}
			log.Debug(fmt.Sprintf("Renamed '%s' into '%s'", oldpath, newpath))
			if strings.HasSuffix(name, ".md") {
				oldTitle := strings.TrimSuffix(name, ".md")
				newTitle := newPrefix + strings.TrimSuffix(strings.TrimPrefix(name, oldPrefix), ".md")
				newTitles <- [2]string{oldTitle, newTitle}
			}
		}
	}
}

// Change all prefixes of subdirectories and files in a given directory
func ChangePrefix(dir, oldPrefix, newPrefix string) error {
	//get titles to change
	newTitles := make(chan [2]string)
	quit := make(chan bool)
	newBacklinks := make(map[string]string)
	go fileexplorer.ApplyToLeavesBeforeRootAsync(dir, ChangeFileAndDirPrefix(oldPrefix, newPrefix, newTitles), quit)
	closed := false
	for {
		select {
		case pair := <-newTitles:
			newBacklinks[pair[0]] = pair[1]
		case <-quit:
			close(newTitles)
			closed = true
			for pair := range newTitles {
				newBacklinks[pair[0]] = pair[1]
			}
		}
		if closed {
			break
		}
	}
	//perform change
	quit = make(chan bool)
	go fileexplorer.ApplyToAllFilesAsync(dir, WrapSwitchBackLinks(newBacklinks, make(map[string]string)), quit)
	<-quit
	return nil
}
