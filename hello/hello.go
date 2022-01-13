package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	name := os.Getenv("GO_USER")
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	if name == "" {
		name = "World"
	}

	for {
		fmt.Printf("Hello, %s!\n", name)
		time.Sleep(1 * time.Second)
	}
}
