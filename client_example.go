package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "grpc_tets/bridge_http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Command-line flags
	serverAddr := flag.String("server", "localhost:50051", "The gRPC server address")
	vendor := flag.String("vendor", "example-vendor", "Vendor name")
	username := flag.String("username", "example-user", "Username")
	password := flag.String("password", "example-password", "Password")
	reqType := flag.String("reqtype", "GET", "Request type (GET, POST, etc.)")
	timeout := flag.Duration("timeout", 10*time.Second, "Request timeout")
	flag.Parse()

	// Set up a connection to the server
	log.Printf("Connecting to gRPC server at %s...", *serverAddr)
	
	conn, err := grpc.Dial(*serverAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	
	log.Println("Connected successfully!")

	// Create a client
	client := pb.NewPgCallClient(conn)

	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Prepare the request
	request := &pb.HitRequest{
		Vendor:   *vendor,
		Username: *username,
		Password: *password,
		ReqType:  *reqType,
	}

	// Make the RPC call
	log.Println("---")
	log.Printf("Sending request:")
	log.Printf("  Vendor:   %s", request.Vendor)
	log.Printf("  Username: %s", request.Username)
	log.Printf("  ReqType:  %s", request.ReqType)
	log.Println("---")
	
	response, err := client.HitPg(ctx, request)
	if err != nil {
		log.Fatalf("Error calling HitPg: %v", err)
	}

	// Display the response
	log.Println("Response received:")
	log.Printf("  Status Code: %d", response.StatusCode)
	
	if response.ErrorMessage != "" {
		log.Printf("  Error Message: %s", response.ErrorMessage)
	}
	
	log.Printf("  Response Body Length: %d bytes", len(response.Response))
	
	if len(response.Headers) > 0 {
		log.Println("  Headers:")
		for key, value := range response.Headers {
			log.Printf("    %s: %s", key, value)
		}
	}
	
	if len(response.Response) > 0 {
		log.Printf("  Response Body: %s", string(response.Response))
	}
	
	log.Println("---")
	log.Println("Request completed successfully!")
}
