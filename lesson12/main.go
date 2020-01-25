package main

import (
	"log"
	"os"

	"github.com/andywow/golang-lessons/lesson12/envdir"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("You should start program: <envdir> <cmd> [<cmdargs>]")
	}
	envDirArgs, err := envdir.ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	rc := envdir.RunCmd(os.Args[2:], envDirArgs)
	os.Exit(rc)
}
