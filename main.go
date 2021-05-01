package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func checkNextArg(args []string, errorMessage string) string {
	if len(args) == 0 {
		log.Fatal(errorMessage)
	}
	return args[0]
}

func main() {
	args := os.Args[1:]
	firstarg := checkNextArg(args, "zlk needs at least one command")
	switch firstarg {
	case "change-prefix":
		dir := checkNextArg(args[1:], "Wrong usage: 'zlk change-prefix <dir> <oldprefix> <newprefix>")
		oldprefix := checkNextArg(args[2:], "Wrong usage: 'zlk change-prefix <dir> <oldprefix> <newprefix>")
		newprefix := checkNextArg(args[3:], "Wrong usage: 'zlk change-prefix <dir> <oldprefix> <newprefix>")
		ChangePrefix(dir, oldprefix, newprefix)
	case "change-tag":
		dir := checkNextArg(args[1:], "Wrong usage: 'zlk change-tag <dir> <oldtag> <newtag")
		oldtag := checkNextArg(args[2:], "Wrong usage: 'zlk change-tag <dir> <oldtag> <newtag>")
		newtag := checkNextArg(args[3:], "Wrong usage: 'zlk change-tag <dir> <oldtag> <newtag>")
		ChangeTag(dir, oldtag, newtag)
	}
}
