package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/client/config"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
)

type updateCommandOptions struct {
	Event   *eventapi.Event
	StrTime string
}

// UpdateCmd update command
func UpdateCmd(opts *config.ClientOptions) *cobra.Command {

	cmdOpts := updateCommandOptions{
		Event: &eventapi.Event{},
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update event command",
		Run: func(cmd *cobra.Command, args []string) {

			if cmdOpts.StrTime != "" {
				eventTime, err := time.Parse(timeFormat, cmdOpts.StrTime)
				if err != nil {
					log.Fatal("could not parse time")
				}
				cmdOpts.Event.StartTime = &eventTime
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			connection, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", opts.GRPCHost, opts.GRPCPort),
				grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			defer connection.Close()
			log.Println("Connected to remote server")

			client := eventapi.NewApiServerClient(connection)

			_, err = client.UpdateEvent(context.Background(), cmdOpts.Event)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("event updated")
		},
	}

	cmd.PersistentFlags().StringVar(&cmdOpts.Event.Uuid, "uuid", "", "event uuid")
	cmd.PersistentFlags().StringVar(&cmdOpts.StrTime, "starttime", "", "event start time - 2018.01.02 20:00")
	cmd.PersistentFlags().Int64Var(&cmdOpts.Event.Duration, "duration", 0, "event duration")
	cmd.PersistentFlags().StringVar(&cmdOpts.Event.Header, "header", "", "event header")
	cmd.PersistentFlags().StringVar(&cmdOpts.Event.Description, "description", "", "event description")

	if err := cmd.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	return cmd

}
