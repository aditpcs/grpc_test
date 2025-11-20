package main

import (
	"context"
	"log"
	"time"

	pb "grpc_tets/bridge_http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server
	serverAddress := "localhost:50051"

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewPgCallClient(conn)

	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Prepare the request
	request := &pb.HitRequest{
		Vendor:   "inacash",
		Username: "susmaili.bungsu@pcsindonesia.co.id",
		Password: "92198316",
		ReqType:  "GET",
	}

	// Make the RPC call
	log.Printf("Calling HitPg with vendor: %s, username: %s, reqType: %s",
		request.Vendor, request.Username, request.ReqType)

	response, err := client.HitPg(ctx, request)
	if err != nil {
		log.Fatalf("Error calling HitPg: %v", err)
	}

	// Display the response
	log.Printf("Response received:")
	log.Printf("  Status Code: %d", response.StatusCode)
	log.Printf("  Error Message: %s", response.ErrorMessage)
	log.Printf("  Response Body Length: %d bytes", len(response.Response))
	log.Printf("  Headers: %v", response.Headers)

	if len(response.Response) > 0 {
		log.Printf("  Response Body: %s", string(response.Response))
	}
}
