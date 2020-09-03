package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	pb "go-kit-grpc-srv/service/mpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
)

//创建bookList的EndPoint
func makeSayHelloMyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if r,ok:=request.(*pb.HelloRequest); ok {
			return &pb.HelloReply{Message: "Hello againMy " + r.GetName()}, nil
		}
		return  nil, nil
	}
}

func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeResponse(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	sayHelloMyHandler kitgrpc.Handler
}

func (s *server) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (s *server) SayHelloMy(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	_, rsp, err := s.sayHelloMyHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*pb.HelloReply), err
}


func main() {
	var grpcAddr string = ":51002"
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	logger.Log("transport", "gRPC", "addr", grpcAddr)

	grpcListener, err := net.Listen("tcp", grpcAddr)

	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	server1 := &server{}
	server1.sayHelloMyHandler = kitgrpc.NewServer(makeSayHelloMyEndpoint(), decodeRequest, encodeResponse)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, server1)
	if err := s.Serve(grpcListener); err != nil {
		logger.Log("errmsg",fmt.Sprintf("failed to serve: %v", err))
	}

}
