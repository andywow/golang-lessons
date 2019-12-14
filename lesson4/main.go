package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/andywow/golang-lessons/lesson4/dictparser"
)

var (
	fileName = flag.String("filename", "", "file to use as input")
	data     = flag.String("data", "", "Data to parse")
)

func main() {
	flag.Parse()
	if *fileName == "" && *data == "" {
		fmt.Fprintf(os.Stderr, "You need specify filename or data parameter")
		os.Exit(1)
	}
	if *fileName != "" && *data != "" {
		fmt.Fprintf(os.Stderr, "You need specify only one of filename and data parameters")
		os.Exit(1)
	}
	result := []string{}
	if *data != "" {
		fmt.Println("Parsing data from input")
		result = dictparser.Top10(*data)
	}
	if *fileName != "" {
		content, err := ioutil.ReadFile(*fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		text := string(content)
		fmt.Println("Parsing data from file")
		result = dictparser.Top10(text)
	}
	fmt.Println("Top10 words:")
	fmt.Println(strings.Join(result, "\n"))
}
