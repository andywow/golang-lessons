package listdatecmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/client/config"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const dateFormat = "2006.01.02"

type options struct {
	EventDate *eventapi.EventDate
	date      string
}

var cmdOpts options

// MakeCmd create command
func MakeCmd(opts *config.ClientOptions) *cobra.Command {

	cmdOpts = options{
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

			connection, err := grpc.Dial(fmt.Sprintf("%s:%d", opts.GRPCHost, opts.GRPCPort),
				grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
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
	cmd.MarkPersistentFlagRequired("date")

	return cmd

}
