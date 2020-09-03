package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	pb "go-kit-grpc-client/service/mpb"
	"google.golang.org/grpc"
	"os"
	"time"
)

const (
	address     = "localhost:51002"
	defaultName = "world"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	logger.Log("transport", "gRPC", "addr", address)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Log("errmsg", fmt.Sprintf("did not connect: %v", err))
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHelloMy(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		logger.Log("errmsg", fmt.Sprintf("could not greet: %v", err))
	}
	logger.Log("msg", fmt.Sprintf("Greeting: %s", r.GetMessage()))
}
