package main

import (
	"log"
	"os"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/factory"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	StorageType string `default:"memory" split_words:"true"`
}

func main() {

	var cfg config

	err := envconfig.Process("calendar", &cfg)
	if err != nil {
		envconfig.Usage("calendar", cfg)
		os.Exit(1)
	}

	//TODO process repository later
	_, err = factory.GetRepository(cfg.StorageType)
	if err != nil {
		//TODO replace with zap
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	log.Println("Finished")
}
