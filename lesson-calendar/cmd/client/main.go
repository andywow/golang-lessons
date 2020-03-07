package main

import (
	"log"

	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/createcmd"
	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/deletecmd"
	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/listdatecmd"
	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/listmonthcmd"
	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/listweekcmd"
	"github.com/andywow/golang-lessons/lesson-calendar/cmd/client/updatecmd"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/client/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(createcmd.MakeCmd(&options))
	rootCmd.AddCommand(updatecmd.MakeCmd(&options))
	rootCmd.AddCommand(deletecmd.MakeCmd(&options))
	rootCmd.AddCommand(listdatecmd.MakeCmd(&options))
	rootCmd.AddCommand(listweekcmd.MakeCmd(&options))
	rootCmd.AddCommand(listmonthcmd.MakeCmd(&options))

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}
