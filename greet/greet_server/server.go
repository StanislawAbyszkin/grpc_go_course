package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/StanislawAbyszkin/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greet_pb.GreetRequest) (*greet_pb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName

	res := &greet_pb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greet_pb.GreetManyTimesRequest, stream greet_pb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number: " + strconv.Itoa(i)
		res := &greet_pb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greet_pb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading the client stream
			res := &greet_pb.LongGreetResponse{
				Result: result,
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatal("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func main() {
	fmt.Println("Starting gRPC server!")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greet_pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
