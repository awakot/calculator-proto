package main

import (
	"fmt"
	"log"
	"time"

	"context"

	"io"

	calculatorproto "github.com/waytkheming/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello from client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect:%v", err)
	}

	defer conn.Close()

	c := calculatorproto.NewCalculatorServiceClient(conn)
	// fmt.Printf("created client: %f", c)

	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	// doErrorUnary(c)

}

func doUnary(c calculatorproto.CalculatorServiceClient) {

	req := &calculatorproto.SumReq{
		FirstNum:  11,
		SecondNum: 22,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GRPC: %v", err)
	}
	log.Printf("Response from Calculator: %v", res.Result)
}

func doClientStreaming(c calculatorproto.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaimng rpc....")

	req := []*calculatorproto.ComputeAvarageReq{
		&calculatorproto.ComputeAvarageReq{
			Num: 1,
		},
		&calculatorproto.ComputeAvarageReq{
			Num: 1,
		},
		&calculatorproto.ComputeAvarageReq{
			Num: 12,
		},
		&calculatorproto.ComputeAvarageReq{
			Num: 18,
		},
		&calculatorproto.ComputeAvarageReq{
			Num: 22,
		},
	}
	stream, err := c.ComputeAvarage(context.Background())
	if err != nil {
		log.Fatalf("error while reading ComputeAvarage: %v", err)
	}

	for _, request := range req {
		fmt.Printf("Sending req: %v\n", request)
		stream.Send(request)
		time.Sleep(1000 * time.Millisecond)
	}

}

func doBiDiStreaming(c calculatorproto.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaimng rpc....")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int64{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorproto.FindMaximumRequest{
				Num: number,
			})
			time.Sleep(1000 * time.Millisecond)
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
				log.Fatalf("Error while Receiving: %v", err)
				break
			}
			maximum := res.GetResult()
			fmt.Printf("Received: %v\n", maximum)
		}

	}()
	<-waitc

}

func doErrorUnary(c calculatorproto.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaimng rpc....")
	num := int32(-10)
	res, err := c.SquareRoot(context.Background(), &calculatorproto.SquareRootRequest{Num: num})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("nagetive num error")
			}
		} else {
			log.Fatalf("Big error: %v", err)
		}
	}
	fmt.Printf("result of root: %v", res.GetNumRoot())
}
