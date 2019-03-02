package main

import (
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

	doSum(c)
}

func doSum(c calculatorpb.CalculatorServiceClient){
	req := calculatorpb.SumRequest{
		A: 10,
		B: 20,
	}

	res, err := c.Sum(context.Background(), &req)
	if err != nil{
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}

	fmt.Printf("Sum of a=%v and b=%v is %v\n", req.A, req.B, res.Result)
}