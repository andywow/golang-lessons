package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatalf("You need specify filename or data parameter")
	}
	if *fileName != "" && *data != "" {
		log.Fatalf("You need specify only one of filename and data parameters")
	}
	var result []string
	if *data != "" {
		log.Println("Parsing data from input")
		result = dictparser.Top10(*data)
	}
	if *fileName != "" {
		content, err := ioutil.ReadFile(*fileName)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		text := string(content)
		log.Println("Parsing data from file")
		result = dictparser.Top10(text)
	}
	log.Println("Top10 words:")
	fmt.Println(strings.Join(result, "\n"))
}
