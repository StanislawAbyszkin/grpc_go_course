package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBiderectionalStreaming(c)
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

func doClientStreaming(c greet_pb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greet_pb.LongGreetRequest{
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Stasiu",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Nati",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Kasia",
			},
		},
		&greet_pb.LongGreetRequest{
			Greeting: &greet_pb.Greeting{
				FirstName: "Mis",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling server with clinet streaming, err: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiveing response from server, err: %v", err)
	}

	fmt.Printf("Long Response: %v\n", res.Result)
}

func doBiderectionalStreaming(c greet_pb.GreetServiceClient){
		fmt.Println("Starting to do a Bi directional Streaming RPC")

		// we create a stream by invoking the client
		stream, err := c.GreetEveryone(context.Background())
		if err != nil {
			log.Fatalf("Error while creating stream; %v\n", err)
			return
		}

		requests := []*greet_pb.GreetEveryoneRequest{
			&greet_pb.GreetEveryoneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Stasiu",
				},
			},
			&greet_pb.GreetEveryoneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Nati",
				},
			},
			&greet_pb.GreetEveryoneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Kasia",
				},
			},
			&greet_pb.GreetEveryoneRequest{
				Greeting: &greet_pb.Greeting{
					FirstName: "Mis",
				},
			},
		}

		waitc := make(chan struct{})
		// we send a bunch of messages to the clinet
		go func() {
			// function to send a bunch of messages
			for _, req := range requests{
				fmt.Printf("Sending message: %v\n", req)
				stream.Send(req)
				time.Sleep(1000*time.Millisecond)
			}
			stream.CloseSend()
		}()

		// we receive a bunch of messages from the client
		go func() {
			// function to receive a bunch of messages
			for {
				res, err := stream.Recv()
				if err == io.EOF{
					break
				}
				if err != nil {
					log.Fatalf("Error while receiving; %v\n", err)
					break
				}
				fmt.Printf("Received: %v", res.GetResult())
			}
			close(waitc)
		}()

		// block until everything is done
		<-waitc
}