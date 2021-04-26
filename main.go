package main

import (
	"fmt"

	"github.com/felixgaschi/smarterzettelkasten/file_explorer"
)

func promptStringAsync(s1, s2 string, ch chan bool) {
	fmt.Printf("%s/%s\n", s1, s2)
	ch <- true
}

func main() {
	ch := make(chan bool)
	go file_explorer.ApplyToLeavesBeforeRootAsync(".", ChangePrefixAsync("file_", "file_"), ch)
	<-ch
}
