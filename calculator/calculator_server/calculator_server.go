package main

import (
	"io"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"github.com/StanislawAbyszkin/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)

	a := req.GetA()
	b := req.GetB()

	res := &calculatorpb.SumResponse{
		Result: a + b,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition rpc call with %v\n", req)
	number := float64(req.GetNumberToDecompose())
	k := 2.
	for number > 1 {
		if math.Mod(number, k) == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Factor: int32(k),
			}
			stream.Send(res)
			number = math.Floor(number / k)
			time.Sleep(200 * time.Millisecond)
		} else {
			k++
		}


	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked with streaming request")

	total := 0.
	cnt := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF{
			res := &calculatorpb.ComputeAverageResponse{
				Average: float64(total / float64(cnt)),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatal("Error while reading client stream: %v", err)
		}
		total += float64(req.GetValue())
		cnt ++
	}
}

func main() {
	fmt.Println("Starting gRPC server!")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
