package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/StanislawAbyszkin/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greet_pb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)

	doServerStreaming(c)
}

func doUnary(c greet_pb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &greet_pb.GreetRequest{
		Greeting: &greet_pb.Greeting{
			FirstName: "Stanislaw",
			LastName:  "Abyszkin",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greet_pb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

	req := &greet_pb.GreetManyTimesRequest{
		Greeting: &greet_pb.Greeting{
			FirstName: "Stanislaw",
			LastName:  "Abyszkin",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			log.Printf("Reached stream EOF")
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}
