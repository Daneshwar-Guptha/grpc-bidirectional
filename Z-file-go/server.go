package main

import (
	"io"
	"log"
	"net"
	"os"

	pb "grpc-go/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

// IMPORTANT: stream type MUST be MessageService_SayHelloServer
func (s *server) SayHello(
	req *pb.RequestMessage,
	stream pb.MessageService_SayHelloServer,
) error {

	file, err := os.Open("./files/" + req.GetName())
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 64*1024) // 64KB chunks (large file safe)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&pb.ResponseMessage{
			Content: buffer[:n],
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, &server{})

	log.Println("Server running on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
