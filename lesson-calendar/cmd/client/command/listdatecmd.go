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

type listDateCommandOptions struct {
	EventDate *eventapi.EventDate
	date      string
}

// ListDateCmd list date command
func ListDateCmd(opts *config.ClientOptions) *cobra.Command {

	cmdOpts := listDateCommandOptions{
		EventDate: &eventapi.EventDate{},
	}

	cmd := &cobra.Command{
		Use:   "listweek",
		Short: "get events for week command",
		Run: func(cmd *cobra.Command, args []string) {

			eventTime, err := time.Parse(dateFormat, cmdOpts.date)
			if err != nil {
				log.Fatal("could not parse time")
			}
			cmdOpts.EventDate.Date = &eventTime

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

			var events *eventapi.EventList
			if events, err = client.GetEventsForWeek(context.Background(), cmdOpts.EventDate); err != nil {
				log.Fatalf("could not get events: %v", err)
			}
			log.Println("event list:")
			for _, remoteEvent := range events.Events {
				log.Printf("uuid: %s, time: %s\n", remoteEvent.Uuid, remoteEvent.StartTime)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&cmdOpts.date, "date", "", "events date - 2019.03.04")
	if err := cmd.MarkPersistentFlagRequired("date"); err != nil {
		log.Fatal(err)
	}

	return cmd

}
