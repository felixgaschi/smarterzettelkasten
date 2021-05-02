package main

import (
	fileexplorer "github.com/felixgaschi/smarterzettelkasten/fileexplorer"
)

// Change all occurence of a tag name in all notes in a directory
func ChangeTag(dir, oldtag, newtag string) {
	tagChange := make(map[string]string)
	tagChange[oldtag] = newtag

	quit := make(chan bool)
	go fileexplorer.ApplyToAllFilesAsync(dir, WrapSwitchBackLinks(make(map[string]string), tagChange), quit)
	<-quit
}
