package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		if err == io.EOF {
			res := &calculatorpb.ComputeAverageResponse{
				Average: float64(total / float64(cnt)),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatal("Error while reading client stream: %v", err)
		}
		total += float64(req.GetValue())
		cnt++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a bi-directional streaming request\n")
	currMax := int64(-999999999) // start with pretty low min int

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream; %v\n", err)
		}

		currMax = max(currMax, req.GetNumber())

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			CurrentMax: currMax,
		})
		if err != nil {
			log.Fatalf("Error while sending to client; %v\n", err)
		}
	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Received SquareRoot RPC\n")

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v", number))
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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

func max(a, b int64) int64 {
	if a >= b {
		return a
	} else {
		return b
	}
}
