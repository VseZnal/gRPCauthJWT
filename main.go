package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"

	"gRPCauthJWT/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	go GrpcServer()
	go GrpcClient()
	var a string
	fmt.Scan(&a)
}

func GrpcServer() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := api.Server{}
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
func GrpcClient() {
	var conn *grpc.ClientConn
	//call Login
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewPingClient(conn)
	loginReply, err := c.Login(context.Background(), &api.LoginRequest{Username: "Slava", Password: "Slava"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	fmt.Println("Login Reply:", loginReply)
	//Call SayHello
	requestToken := new(api.AuthToekn)
	requestToken.Token = loginReply.Token
	conn, err = grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(requestToken))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c = api.NewPingClient(conn)
	helloreply, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", helloreply.Greeting)
}
