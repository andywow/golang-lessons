package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/logconfig"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository/dbstorage"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/grpc/apiserver"
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
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("could not sync logger: %v\n", err)
		}
	}()

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

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.HTTPListen, nil); err != nil {
			sugar.Fatal("cannot start /metrics endpoint: ", err)
		}
	}()

	if err := apiServer.StartServer(ctx, cfg.GRPCListen,
		apiserver.WithLogger(logger), apiserver.WithRepository(&repository)); err != nil {
		sugar.Errorf("failed start api server: %v", err)
		os.Exit(1)
	}

}
