package main

func SendTrue(ch chan bool) {
	ch <- true
}
