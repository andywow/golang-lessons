package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/msgsystem/rabbitmq"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/sender"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/logconfig"
)

var cfg config.Config

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

	rabbitmq, err := rabbitmq.NewRabbitMQ(ctx, cfg.RabbitMQ)
	if err != nil {
		sugar.Fatalf("error, while connecting to message system: %s\n", err)
	}
	sugar.Info("Message system initialized")

	cron := sender.Sender{}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		sugar.Infof("received signal: %s", sig)
		cancel()
	}()

	cron.Start(ctx,
		sender.WithLogger(logger), sender.WithMsgSystem(&rabbitmq))

}
