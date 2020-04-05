package main

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
)

var apiServerAddress = os.Getenv("APISERVER_ADDRESS")

func init() {
	if apiServerAddress == "" {
		apiServerAddress = "127.0.0.1:9090"
	}
}

type apiServerTestClient struct {
	cancelCall, cancelConn context.CancelFunc
	connection             *grpc.ClientConn
	client                 eventapi.ApiServerClient
	callCtx                context.Context
}

func (t *apiServerTestClient) create() {

	var (
		ctx context.Context
		err error
	)
	ctx, t.cancelConn = context.WithTimeout(context.Background(), 5*time.Second)

	if t.connection, err = grpc.DialContext(ctx, apiServerAddress,
		grpc.WithInsecure(), grpc.WithBlock()); err != nil {
		panic(err)
	}

	t.client = eventapi.NewApiServerClient(t.connection)

	t.callCtx, t.cancelCall = context.WithTimeout(context.Background(), 5*time.Second)

}

func (t *apiServerTestClient) close() {
	t.cancelCall()
	t.cancelConn()
	if err := t.connection.Close(); err != nil {
		panic(err)
	}
}
