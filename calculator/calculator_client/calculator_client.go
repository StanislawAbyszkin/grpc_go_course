package main

import (
	"time"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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

	// doPrimeNumberDecomposition(c)

	// doAverageSum(c)

	doFindMaximum(c)
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
		if err == io.EOF {
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

func doAverageSum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	scanner := bufio.NewScanner(os.Stdin)

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling server with clinet streaming, err: %v", err)
	}

	fmt.Println("Input values")
	fmt.Println("---------------------")
	fmt.Print("-> ")
	for scanner.Scan() {
		text := scanner.Text()
		// text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("", text) == 0 {
			fmt.Println("Done sending values")
			break
		}

		value, err := strconv.ParseInt(text, 0, 64)
		if err != nil {
			log.Printf("Couldn't convert given value %v to int, sending existing data to get average", text)
			break
		}

		req := &calculatorpb.ComputeAverageRequest{
			Value: value,
		}
		stream.Send(req)
		fmt.Print("-> ")
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiveing response from server, err: %v", err)
	}

	fmt.Printf("Total average: %v\n", res.Average)
}

func doFindMaximum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Bi directional Streaming RPC")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream; %v\n", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	waitc := make(chan struct{})

	fmt.Println("Input values")
	fmt.Println("---------------------")
	fmt.Print("-> ")

	go func() {
		for scanner.Scan() {
			
			text := scanner.Text()
			// text = strings.Replace(text, "\n", "", -1)

			if strings.Compare("", text) == 0 {
				fmt.Println("Done sending values")
				break
			}

			value, err := strconv.ParseInt(text, 0, 64)
			if err != nil {
				log.Printf("Couldn't convert given value %v to int", text)
				break
			}

			req := &calculatorpb.FindMaximumRequest{
				Number: value,
			}
			stream.Send(req)
			time.Sleep(1* time.Millisecond)
			fmt.Print("\n-> ")
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving from server; %v\n", err)
			}
			fmt.Printf("Current max value: %v", res.GetCurrentMax())
		}
		close(waitc)
	}()
	<-waitc
}
