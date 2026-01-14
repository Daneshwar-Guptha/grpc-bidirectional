package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "grpc-go/proto"
	"google.golang.org/grpc"
)

func main() {
	var fileName string
	fmt.Print("Enter file name: ")
	fmt.Scan(&fileName)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	req := &pb.RequestMessage{
		Name: fileName,
	}

	stream, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create("downloaded_" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	var total int64

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Download completed:", total, "bytes")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		n, _ := outFile.Write(resp.Content)
		total += int64(n)
		fmt.Println("Downloaded:", total)
	}
}
