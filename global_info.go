package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/felixgaschi/smarterzettelkasten/file_explorer"
	log "github.com/sirupsen/logrus"
)

type GlobalInfo struct {
	// maps note title to titles of notes that mention it
	backlinks map[string][]string

	// maps tags to titles of notes that mention it
	tagsBacklinks map[string][]string

	// maps not titles to base directory
	titleToBaseDir map[string]string
}

func SendToGlobal(backlinksChannel chan [2]string, tagsChannel chan [2]string, baseChannel chan [2]string) func(string, string, chan bool) {
	return func(base, fname string, quit chan bool) {
		path := path.Join(base, fname)
		title := strings.TrimSuffix(fname, ".md")
		if title == fname {
			quit <- true
			return
		}
		baseChannel <- [2]string{title, path}
		go SendBacklinksAndTags(path, title, quit, backlinksChannel, tagsChannel)
	}
}

func AddBacklink(g GlobalInfo, key, value string) GlobalInfo {
	log.Debug(fmt.Sprintf("Found backlink to '%s' in '%s'", key, value))
	_, exists := g.backlinks[key]
	if !exists {
		g.backlinks[key] = make([]string, 0)
	}
	g.backlinks[key] = append(g.backlinks[key], value)
	return g
}

func AddTagBacklink(g GlobalInfo, key, value string) GlobalInfo {
	log.Debug(fmt.Sprintf("Found tag '%s' in '%s'", key, value))
	_, exists := g.tagsBacklinks[key]
	if !exists {
		g.tagsBacklinks[key] = make([]string, 0)
	}
	g.tagsBacklinks[key] = append(g.tagsBacklinks[key], value)
	return g
}

func AddBaseDir(g GlobalInfo, key, value string) GlobalInfo {
	log.Debug(fmt.Sprintf("Found note '%s' with dir '%s'", key, value))
	g.titleToBaseDir[key] = value
	return g
}

func RefreshGlobalInfo(dir string) GlobalInfo {
	res := GlobalInfo{make(map[string][]string), make(map[string][]string), make(map[string]string)}

	backlinksChannel := make(chan [2]string)
	tagsChannel := make(chan [2]string)
	baseChannel := make(chan [2]string)
	quit := make(chan bool)
	go file_explorer.ApplyToAllFilesAsync(dir, SendToGlobal(backlinksChannel, tagsChannel, baseChannel), quit)
	for {
		closing := false
		select {
		case elem := <-backlinksChannel:
			res = AddBacklink(res, elem[0], elem[1])
		case elem := <-tagsChannel:
			res = AddTagBacklink(res, elem[0], elem[1])
		case elem := <-baseChannel:
			res = AddBaseDir(res, elem[0], elem[1])
		case <-quit:
			close(baseChannel)
			close(tagsChannel)
			close(backlinksChannel)
			closing = true
		}
		if closing {
			for elem := range backlinksChannel {
				res = AddBacklink(res, elem[0], elem[1])
			}
			for elem := range tagsChannel {
				res = AddTagBacklink(res, elem[0], elem[1])
			}
			for elem := range baseChannel {
				res = AddBaseDir(res, elem[0], elem[1])
			}
			break
		}
	}

	return res
}
