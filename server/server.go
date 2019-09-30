package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"io"

	"math"

	calculatorproto "github.com/waytkheming/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorproto.SumReq) (*calculatorproto.SumRes, error) {
	fmt.Printf("calculator function was invoked with %v", req)
	firstNum := req.FirstNum
	secondNum := req.SecondNum
	sum := firstNum + secondNum
	res := &calculatorproto.SumRes{
		Result: sum,
	}

	return res, nil
}

func (*server) ComputeAvarage(stream calculatorproto.CalculatorService_ComputeAvarageServer) error {
	fmt.Printf("ComputeAvarage function was invoked with a streaming request \n")
	i := int64(0)
	sum := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result := float64(sum) / float64(i)
			return stream.SendAndClose(&calculatorproto.ComputeAvarageRes{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while: %v", err)
		}
		num := req.GetNum()
		i++
		sum = sum + num
		fmt.Println(i)
		fmt.Println(sum)
	}
}

func (*server) FindMaximum(stream calculatorproto.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request \n")
	maximum := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {

			return nil
		}
		if err != nil {
			log.Fatalf("error while: %v", err)
			return err
		}
		number := req.GetNum()
		if number > maximum {
			maximum = number
			sendErr := stream.Send(&calculatorproto.FindMaximumResponse{
				Result: maximum,
			})
			if sendErr != nil {
				log.Fatalf("error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorproto.SquareRootRequest) (*calculatorproto.SquareRootResponse, error) {
	fmt.Printf("Squareroot function was invoked with a streaming request \n")

	number := req.GetNum()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number: %v", number),
		)
	}
	return &calculatorproto.SquareRootResponse{
		NumRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorproto.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
