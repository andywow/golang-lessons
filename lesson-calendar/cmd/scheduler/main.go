package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/msgsystem/rabbitmq"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/logconfig"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository/dbstorage"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/scheduler"
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
		sugar.Fatalf("error, while connecting to database: %s\n", err)
	}
	sugar.Info("Storage initialized")

	rabbitmq, err := rabbitmq.NewRabbitMQ(ctx, cfg.RabbitMQ)
	if err != nil {
		sugar.Fatalf("error, while connecting to message system: %s\n", err)
	}
	sugar.Info("Message system initialized")

	cron := scheduler.Scheduler{}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		sugar.Infof("received signal: %s", sig)
		sugar.Info("wait maximum 1 minute for correct termination")
		cancel()
	}()

	cron.Start(ctx,
		scheduler.WithLogger(logger), scheduler.WithRepository(&repository), scheduler.WithMsgSystem(&rabbitmq))

}
