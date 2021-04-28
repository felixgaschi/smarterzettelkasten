package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

const (
	normal       = iota
	open_bracket = iota
	in_tag       = iota
	in_link      = iota
)

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
