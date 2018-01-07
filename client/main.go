/*
gRPC Client
*/

package main

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"

	nso "github.com/nleiva/slack-nso/nso"
	pb "github.com/nleiva/slack-nso/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "grpc.nleiva.com:50051"
)

// GetCmd subscribes to a stream of commands, returns a channel.
func GetCmd(client pb.CommSvcClient) chan []string {
	// 's' is the string channel where the data received will be sent.
	s := make(chan []string)
	stream, err := client.GetCmd(context.Background(), &pb.Id{})
	if err != nil {
		log.Fatalf("Server says: %v", err)
	}
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Server says: %v", err)
			}
			s <- []string{res.GetCmd(), res.GetArg1(), res.GetArg2()}
		}
		close(s)
	}()
	return s
}

func main() {
	api := new(url.URL)
	api.Scheme = "http"
	api.Host = os.Getenv("NSO_SERVER")
	api.User = url.UserPassword(os.Getenv("NSO_USER"), os.Getenv("NSO_PASSWORD"))
	device := os.Getenv("NSO_DEVICE")
	s := new(nso.Server)
	s.Addr = api

	// Security options
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	// Set up a secure connection to the server.
	conn, err := grpc.Dial(address, opts...)
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCommSvcClient(conn)
	ch := GetCmd(client)

	for msg := range ch {
		s.StaticRoute(msg, device)
	}
}
