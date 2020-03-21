package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository/dbstorage"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/grpc/apiserver"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/logconfig"
)

func main() {

	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger, err := logconfig.GetLoggerForConfig(cfg)
	if err != nil {
		fmt.Printf("could not configure logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repository, err := dbstorage.NewDatabaseStorage(ctx, cfg.DB)
	if err != nil {
		sugar.Fatalf("failed initialize storage: %v", err)
	}
	logger.Info("Storage initialized")

	sugar.Infof("Starting server on %s", cfg.GRPCListen)

	apiServer := apiserver.APIServer{}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		sugar.Infof("received signal: %s", sig)
		cancel()
	}()

	apiServer.StartServer(ctx, cfg.GRPCListen,
		apiserver.WithLogger(logger), apiserver.WithRepository(&repository))

}
