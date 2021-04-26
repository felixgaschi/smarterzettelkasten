package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func ChangePrefixAsync(oldPrefix, newPrefix string) func(string, string, chan bool) {
	return func(base string, fname string, quit chan bool) {
		if strings.HasPrefix(fname, oldPrefix) {
			newFname := newPrefix + strings.TrimPrefix(fname, oldPrefix)
			err := os.Rename(path.Join(base, fname), path.Join(base, newFname))
			if err != nil {
				fmt.Println(err)
			}
		}
		quit <- true
	}
}
