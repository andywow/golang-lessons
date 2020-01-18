package main

import (
	"flag"
	"log"

	"github.com/andywow/golang-lessons/lesson10/copyfile"
)

var (
	fromFileName   = flag.String("from", "", "source file")
	toFileName     = flag.String("to", "", "destionation file")
	fromFileOffset = flag.Int("offset", 0, "source file offset")
	fromFileBytes  = flag.Int("limit", 0, "bytes to copy")
)

func main() {
	flag.Parse()
	if *fromFileName == "" {
		log.Fatalf("You need specify source filename")
	}
	if *toFileName == "" {
		log.Fatalf("You need specify destionation filename")
	}
	//bar := pb.New(*fromFileBytes)
	err := copyfile.Copy(*fromFileName, *toFileName, *fromFileBytes, *fromFileOffset)
	if err != nil {
		log.Fatalf("Error, while copy file: %s", err)
	}
}
