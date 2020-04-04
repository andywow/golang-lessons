package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/command"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/client/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "client",
		Short: "client for grpc api server",
		Long:  "client for testing grpc api server",
	}
)

func init() {

	options := config.ClientOptions{}

	rootCmd.PersistentFlags().StringVar(&options.GRPCHost, "host", "127.0.0.1", "host of grpc server")
	rootCmd.PersistentFlags().Int64Var(&options.GRPCPort, "port", 9090, "port of grpc server")
	if err := viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port")); err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(command.CreateCmd(&options))
	rootCmd.AddCommand(command.UpdateCmd(&options))
	rootCmd.AddCommand(command.DeleteCmd(&options))
	rootCmd.AddCommand(command.ListDateCmd(&options))
	rootCmd.AddCommand(command.ListWeekCmd(&options))
	rootCmd.AddCommand(command.ListMonthCmd(&options))

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}
