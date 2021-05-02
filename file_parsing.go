package main

import (
	"fmt"
	"io/ioutil"
	"path"

	log "github.com/sirupsen/logrus"
)

const (
	normal       = iota
	open_bracket = iota
	in_tag       = iota
	in_link      = iota
)

// parse a file and send backlinks and tags to channels
func SendBacklinksAndTags(fpath, title string, quit chan bool, backlinksChannel chan [2]string, tagsChannel chan [2]string) {
	defer SendTrue(quit)

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Error(fmt.Sprintf("Could not read %s", fpath))
		return
	}

	state := normal
	name := make([]byte, 0)
	for _, b := range data {
		switch state {
		case normal:
			if b == '[' {
				state = open_bracket
			}
		case open_bracket:
			switch b {
			case '[':
				state = in_link
			case '#':
				state = in_tag
			default:
				state = normal
			}
		case in_tag:
			switch b {
			case ']':
				state = normal
				tagsChannel <- [2]string{string(name), string(title)}
				name = make([]byte, 0)
			default:
				name = append(name, b)
			}
		case in_link:
			switch b {
			case ']':
				state = normal
				backlinksChannel <- [2]string{string(name), string(title)}
				name = make([]byte, 0)
			default:
				name = append(name, b)
			}
		}
	}
}

// closure for parametrizing SitchBackLinks
func WrapSwitchBackLinks(newBacklinks map[string]string, newTags map[string]string) func(string, string, chan bool) {
	return func(base, fname string, quit chan bool) {
		fpath := path.Join(base, fname)
		SwitchBackLinks(fpath, newBacklinks, newTags, quit)
	}
}

// parse a file and replace note backlink and tags according to
// maps passed in arguments
func SwitchBackLinks(fpath string, newBacklinks map[string]string, newTags map[string]string, quit chan bool) {
	defer SendTrue(quit)
	//TODO add lock to verify we are not already editing this file
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Error(fmt.Sprintf("Could not read %s", fpath))
		return
	}

	state := normal
	name := make([]byte, 0)
	new_bytes := make([]byte, 0, len(data))
	for _, b := range data {
		switch state {
		case normal:
			new_bytes = append(new_bytes, b)
			if b == '[' {
				state = open_bracket
			}
		case open_bracket:
			new_bytes = append(new_bytes, b)
			switch b {
			case '[':
				state = in_link
			case '#':
				state = in_tag
			default:
				state = normal
			}
		case in_tag:
			switch b {
			case ']':
				state = normal
				new, exists := newTags[string(name)]
				var new_name_bytes []byte
				if exists {
					new_name_bytes = []byte(new)
				} else {
					new_name_bytes = name
				}
				new_bytes = append(new_bytes, new_name_bytes...)
				name = make([]byte, 0)
				new_bytes = append(new_bytes, b)
			default:
				name = append(name, b)
			}
		case in_link:
			switch b {
			case ']':
				state = normal
				new, exists := newBacklinks[string(name)]
				var new_name_bytes []byte
				if exists {
					new_name_bytes = []byte(new)
				} else {
					new_name_bytes = name
				}
				new_bytes = append(new_bytes, new_name_bytes...)
				name = make([]byte, 0)
				new_bytes = append(new_bytes, b)
			default:
				name = append(name, b)
			}
		}
	}

	write_err := ioutil.WriteFile(fpath, new_bytes, 0644)
	if write_err != nil {
		log.Fatal(fmt.Sprintf("Failed to write file '%s'", fpath))
		panic(write_err)
	}
}
