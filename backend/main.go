package main

import "fmt"

type Message struct {
	MsgType string `json:"type"`
}

func main() {
	fmt.Println("Hello Nix!")
}
