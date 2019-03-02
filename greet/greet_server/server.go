package main

import (
	"context"
	"github.com/StanislawAbyszkin/grpc-go-course/greet/greetpb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greet_pb.GreetRequest) (*greet_pb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName

	res := &greet_pb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Starting gRPC server!")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greet_pb.RegisterGreetServiceServer(s, &server{})

	if err:= s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
