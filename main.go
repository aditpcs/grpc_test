package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "grpc_tets/bridge_http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define command-line flags
	requestType := flag.String("type", "", "The type of request to make (e.g., inacash-generate, bri-generate, bri-check-status)")
	referenceNo := flag.String("reference_no", "", "The reference number for check status")
	flag.Parse()

	// Set up a connection to the server
	serverAddress := "localhost:50051"
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewPgCallClient(conn)

	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request *pb.HitRequest

	// Prepare the request based on the flag
	switch *requestType {
	case "inacash-generate":
		request = &pb.HitRequest{
			Vendor:   "inacash",
			Username: "susmaili.bungsu@pcsindonesia.co.id",
			Password: "92198316",
			ReqType:  "create",
			Request: []byte(`{
				"product_code": "QRIS_DIRECT",
				"amount": "1",
				"remark": "payment pose",
				"client_reff": "sample10030901",
				"merchant_id": "INA-B7417383445"
				}`,
			),
			TokenFcm: "eb9eSlxIQaeQhtp2vbp2Vz:APA91bHs3MHpRQMVOoTjSmdq5qLOhpGlQyMeuDxNcF-FtC24_WXmFqizpRZ5zQKkcvO0sLDQ6mT-CAxy3Sc50wUTH4x3cBifjD8WScEt_Pyh8yCYEBKcnrY",
		}
		// Make the RPC call
		log.Printf("Calling HitPg with vendor: %s, reqType: %s",
			request.Vendor, request.ReqType)

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
	case "bri-generate":
		clientReff := time.Now().Unix()
		// Initialize your byte slice (optionally with a pre-defined capacity for speed)
		var req []byte

		// Use Appendf to build the JSON directly into the slice
		req = fmt.Appendf(req,
			`{"amount":{"currency":"IDR","value":"1.00"},"partnerReferenceNo":"%d","merchantId":"000001019000014","terminalId":"10049259"}`,
			clientReff,
		)
		request = &pb.HitRequest{
			Vendor:  "bri",
			ReqType: "generate",
			Request: req,
		}
		// Make the RPC call
		log.Printf("Calling HitPg with vendor: %s, reqType: %s",
			request.Vendor, request.ReqType)

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

		// Unmarshal the response to get the referenceNo
		var respData map[string]any
		if err := json.Unmarshal(response.Response, &respData); err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}

		referenceNo, ok := respData["referenceNo"].(string)
		if !ok {
			log.Fatalf("Could not find referenceNo in response body")
		}

		// Store the reference number
		err = os.WriteFile("/tmp/referenceno", []byte(referenceNo), 0o644)
		if err != nil {
			log.Fatalf("Failed to write reference number to file: %v", err)
		}

	case "bri-check-status":
		var refNo string
		if *referenceNo != "" {
			refNo = *referenceNo
		} else {
			data, err := os.ReadFile("/tmp/referenceno")
			if err != nil {
				log.Fatalf("Reference number not provided and could not be read from file: %v", err)
			}
			refNo = strings.TrimSpace(string(data))
		}

		var req []byte
		req = fmt.Appendf(req,
			`{"serviceCode":"17","originalReferenceNo":"%s", "additionalInfo":{"terminalId":"10049259"}}`,
			refNo,
		)

		request = &pb.HitRequest{
			Vendor:  "bri",
			ReqType: "status",
			Request: req,
		}
		// Make the RPC call
		log.Printf("Calling HitPg with vendor: %s, reqType: %s",
			request.Vendor, request.ReqType)

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
	default:
		fmt.Println("Invalid flag. Use -type=[inacash-generate|bri-generate|bri-check-status]")
		os.Exit(1)
	}
}
