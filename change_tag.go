package main

import (
	file_explorer "github.com/felixgaschi/smarterzettelkasten/file_explorer"
)

// Change all occurence of a tag name in all notes in a directory
func ChangeTag(dir, oldtag, newtag string) {
	tagChange := make(map[string]string)
	tagChange[oldtag] = newtag

	quit := make(chan bool)
	go file_explorer.ApplyToAllFilesAsync(dir, WrapSwitchBackLinks(make(map[string]string), tagChange), quit)
	<-quit
}
