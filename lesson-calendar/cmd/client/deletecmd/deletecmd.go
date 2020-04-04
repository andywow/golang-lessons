package deletecmd

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

var deleteMessage eventapi.EventDelete

// MakeCmd create command
func MakeCmd(opts *config.ClientOptions) *cobra.Command {

	deleteMessage = eventapi.EventDelete{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete event command",
		Run: func(cmd *cobra.Command, args []string) {

			connection, err := grpc.Dial(fmt.Sprintf("%s:%d", opts.GRPCHost, opts.GRPCPort),
				grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			defer connection.Close()
			log.Println("Connected to remote server")

			client := eventapi.NewApiServerClient(connection)
			if _, err := client.DeleteEvent(context.Background(), &deleteMessage); err != nil {
				log.Fatal(err)
			}
			log.Println("event deleted")
		},
	}

	cmd.PersistentFlags().StringVar(&deleteMessage.Uuid, "uuid", "", "event uuid")
	if err := cmd.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	return cmd

}
