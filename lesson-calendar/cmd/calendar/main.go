package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/logconfig"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository/localcache"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/httpserver"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfg config.Config

func init() {
	flag.String("configfile", "config.yaml", "config file path")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// default values
	cfg = config.Config{
		HTTPListen:  "127.0.0.1:8080",
		LogLevel:    "info",
		LogStdout:   true,
		StorageType: "Memory",
	}
}

func main() {

	viper.SetConfigFile(viper.GetString("configfile"))
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("could not read config: %s\n", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf("could not read config: %s\n", err)
		os.Exit(1)
	}

	logger, err := logconfig.GetLoggerForConfig(&cfg)
	if err != nil {
		fmt.Printf("could not configure logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	// cause it was in task
	localcache.NewEventLocalStorage()
	sugar.Info("Storage initialized")

	sugar.Infof("Starting server on %s", cfg.HTTPListen)

	httpserver.StartServer(cfg.HTTPListen, httpserver.WithLogger(logger))

}
