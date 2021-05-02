package main_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	main "github.com/felixgaschi/smarterzettelkasten"
	log "github.com/sirupsen/logrus"
)

const testdir = "./test/test_dir"

func sendBacklinksAndTagsAndWait(fpath, title string) ([][2]string, [][2]string) {
	quit := make(chan bool)
	backlinksChannel := make(chan [2]string)
	tagsChannel := make(chan [2]string)
	go main.SendBacklinksAndTags(fpath, title, quit, backlinksChannel, tagsChannel)
	backlinks := make([][2]string, 0)
	tags := make([][2]string, 0)
	closing := false
	for {
		select {
		case elem := <-backlinksChannel:
			backlinks = append(backlinks, elem)
		case elem := <-tagsChannel:
			tags = append(tags, elem)
		case <-quit:
			close(tagsChannel)
			close(backlinksChannel)
			close(quit)
			closing = true
		}
		if closing {
			break
		}
	}
	for elem := range tagsChannel {
		tags = append(tags, elem)
	}
	for elem := range backlinksChannel {
		backlinks = append(backlinks, elem)
	}
	return backlinks, tags
}

func TestChangeTag(t *testing.T) {
	const title = "O_0.0.1 First note"
	fpath := path.Join(testdir, title+".md")

	path, err := os.Getwd()
	if err != nil {
		t.Log("Could not find current directory")
	} else {
		t.Log(fmt.Sprintf("Current directory: %s", path))
	}

	_, tags := sendBacklinksAndTagsAndWait(fpath, title)

	var tag string
	if len(tags) == 0 {
		t.Error("Failed TestChangeTag because no tag found in " + title)
	} else {
		tag = tags[0][0]
		t.Log("Found tag", tag)
	}

	main.ChangeTag(testdir, tag, tag+tag)

	_, tags = sendBacklinksAndTagsAndWait(fpath, title)

	found := false
	foundtags := make([]string, 0)
	for _, elem := range tags {
		if tag+tag == elem[0] {
			found = true
			break
		}
		foundtags = append(foundtags, elem[0])
	}

	if !found {
		t.Error("Failed TestChangeTage because changed tag not found, instead found: ", foundtags)
	} else {
		t.Log("Found changed tag", tag+tag)
	}

	main.ChangeTag(testdir, tag+tag, tag)
}

func checkElemInDir(path, elem string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("Directory not found", path)
	}
	for _, file := range files {
		if file.Name() == elem {
			return true
		}
	}
	return false
}

func TestChangePrefix(t *testing.T) {
	old_prefix := "N_0"
	new_prefix := "N_1"

	prefixedDir := "N_0.0.1 Test subdir"
	prefixedTitle := "N_0.0.1.1 Note in subdir"
	prefixedFile := prefixedTitle + ".md"

	if !checkElemInDir(testdir, prefixedDir) {
		t.Error(fmt.Sprintf("Before change: could not find '%s' in '%s'", prefixedDir, testdir))
	}
	if !checkElemInDir(path.Join(testdir, prefixedDir), prefixedFile) {
		t.Error(fmt.Sprintf("Before change: could not find '%s' in '%s'", prefixedFile, path.Join(testdir, prefixedDir)))
	}

	main.ChangePrefix(testdir, old_prefix, new_prefix)

	newPrefixedDir := "N_1.0.1 Test subdir"
	newPrefixedTitle := "N_1.0.1.1 Note in subdir"
	newPrefixedFile := newPrefixedTitle + ".md"

	if !checkElemInDir(testdir, newPrefixedDir) {
		t.Error(fmt.Sprintf("After change: could not find '%s' in '%s'", newPrefixedDir, testdir))
	}
	if !checkElemInDir(path.Join(testdir, newPrefixedDir), newPrefixedFile) {
		t.Error(fmt.Sprintf("After change: could not find '%s' in '%s'", newPrefixedFile, path.Join(testdir, newPrefixedDir)))
	}

	main.ChangePrefix(testdir, new_prefix, old_prefix)

	if !checkElemInDir(testdir, prefixedDir) {
		t.Error(fmt.Sprintf("Could not put back '%s' in '%s'", prefixedDir, testdir))
	}
	if !checkElemInDir(path.Join(testdir, prefixedDir), prefixedFile) {
		t.Error(fmt.Sprintf("Could not put back '%s' in '%s'", prefixedFile, path.Join(testdir, prefixedDir)))
	}

}
