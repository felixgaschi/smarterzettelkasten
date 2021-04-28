package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Error(fmt.Sprintf("Smarterzettelkasten requires one and only one argument, got %d", len(args)))
		return
	}
	RefreshGlobalInfo(args[0])
}
