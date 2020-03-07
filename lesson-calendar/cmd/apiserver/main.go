package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository/dbstorage"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/grpc/apiserver"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/logconfig"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository/localcache"
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
		GRPCListen:  "127.0.0.1:9090",
		LogLevel:    "info",
		LogStdout:   true,
		StorageType: "memory",
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

	var repository repository.EventRepository
	switch cfg.StorageType {
	case "database":
		repository, err = dbstorage.NewDatabaseStorage(
			cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword)
		if err != nil {
			sugar.Fatalf("error, while connecting to database: %s\n", err)
		}
	default:
		repository = localcache.NewEventLocalStorage()
	}
	logger.Info("Storage initialized")

	sugar.Infof("Starting server on %s", cfg.GRPCListen)

	apiServer := apiserver.APIServer{}

	apiServer.StartServer(cfg.GRPCListen,
		apiserver.WithLogger(logger), apiserver.WithRepository(&repository))

}
