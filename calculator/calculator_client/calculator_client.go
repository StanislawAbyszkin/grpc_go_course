package main

import (
	"io"
	"context"
	"fmt"
	"log"

	"github.com/StanislawAbyszkin/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cound not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doSum(c)

	doPrimeNumberDecomposition(c)
}

func doSum(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		A: 10,
		B: 20,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}

	fmt.Printf("Sum of a=%v and b=%v is %v\n", req.A, req.B, res.Result)
}

func doPrimeNumberDecomposition(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		NumberToDecompose: 1019494912,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	fmt.Printf("Starting decomposing number: %v\n", req.NumberToDecompose)
	for {
		msg, err := resStream.Recv()
		if err == io.EOF{
			// we've reached the end of the stream
			log.Printf("Reached stream EOF")
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Factor: %v", msg.GetFactor())
	}
}
