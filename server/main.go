/*
gRPC Server
*/

package main

import (
	"log"
	"net"

	pb "github.com/nleiva/slack-nso/proto"
	sl "github.com/nleiva/slack-nso/slack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

// server is used to implement CommSvcServer
type server struct {
	ch chan []string
}

func (s *server) GetCmd(in *pb.Id, stream pb.CommSvc_GetCmdServer) error {
	for r := range s.ch {
		stream.Send(&pb.Command{Cmd: "route", Arg1: r[0], Arg2: r[1]})
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	c := sl.Listen()
	svr := server{ch: c}

	// Security options
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	// Setup a secure Server
	s := grpc.NewServer(opts...)

	pb.RegisterCommSvcServer(s, &svr)
	log.Println("Starting server on port " + port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
