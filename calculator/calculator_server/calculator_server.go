package main

import (
	"context"
	"fmt"
	"log"
	"net"
	
	"github.com/StanislawAbyszkin/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error){
	fmt.Printf("Sum function was invoked with %v\n", req)

	a := req.GetA()
	b := req.GetB()
	
	res := &calculatorpb.SumResponse{
		Result: a + b,
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

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err:= s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
