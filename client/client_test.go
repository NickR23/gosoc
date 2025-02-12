package client_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/NickR23/gosoc/client"
)

// Example test that runs against the server
func TestClientHandshake(t *testing.T) {
	http_client := client.Initialize(&http.Client{})
	handshake_req, _ := client.BuildHandshake()
	resp, err := http_client.Do(handshake_req)
	if err != nil {
		log.Fatal("Got error:", err)
		t.Fatalf("Got error: %v", err)
	}

	if resp.StatusCode != 101 {
		t.Fatalf("Expected 101, got %v", resp.StatusCode)
	}
	log.Println("Response: ", resp.Status)
	log.Println("Response: ", resp.Header)

}
