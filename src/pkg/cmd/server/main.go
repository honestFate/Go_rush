package main

import (
	"log"
	"net"
	api "randsig/pkg/api"
	randsig "randsig/pkg/randomSignaler"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	service := &randsig.GRPCServer{}
	api.RegisterRandomSignalerServer(server, service)

	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Serve(listner); err != nil {
		log.Fatal(err)
	}
}
