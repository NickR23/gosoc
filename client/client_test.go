package client_test

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/NickR23/gosoc/client"
)

// Example test that runs against the server
func TestClientHandshake(t *testing.T) {
	conn, err := client.Handshake("ws://127.0.0.1:9001/")
	if err != nil {
		t.Fatalf("Error during handshake: %v", err)
	}

	message := []byte("Hello world!!!!")

	f := client.WSFrame{
		Fin:        true,
		Opcode:     0x1,
		Mask:       true,
		PayloadLen: uint64(len(message)),
		Payload:    message,
	}

	encodedFrame, _ := f.Encode()
	log.Printf("Encoded Frame: %v", base64.StdEncoding.EncodeToString(encodedFrame))
	_, err = conn.Write(encodedFrame)
	if err != nil {
		t.Fatalf("Error sending frame: %v", err)
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("Error during response read: %v", err)
	}
	log.Printf("Received response: %v", response[:n])
}
