package main

import (
	"fmt"
	"os"

	"github.com/andywow/golang-lessons/lesson2/stringunpack"
)

func main() {
	if (len(os.Args)) == 1 {
		fmt.Fprintf(os.Stderr, "Error: specify input string\n")
		os.Exit(1)
	}
	result, err := stringunpack.Unpack(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Unpacked string: %s\n", result)
}
