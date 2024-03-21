package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/minpeter/grpc-web-example/gen"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9090, "The server port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// SayRepeatHello implements helloworld.GreeterServer
func (s *server) SayRepeatHello(in *pb.RepeatHelloRequest, stream pb.Greeter_SayRepeatHelloServer) error {
	log.Printf("Received: %v", in.GetName())
	for i := range 4 {
		if err := stream.Send(&pb.HelloReply{Message: "Hello " + in.GetName() + fmt.Sprintf(" %d", i)}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// http server

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
