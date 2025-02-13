package client_test

import (
	"encoding/base64"
	"log"
	"net"
	"testing"

	"github.com/NickR23/gosoc/client"
)

func sendFrame(conn net.Conn, f client.WSFrame) ([]byte, error) {
	encodedFrame, _ := f.Encode()
	log.Printf("Encoded Frame: %v", base64.StdEncoding.EncodeToString(encodedFrame))
	_, err := conn.Write(encodedFrame)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}
	// TODO Deserialize this response brother
	log.Printf("Received response: %v", response[:n])
	return response, nil
}

// Example test that runs against the server
func TestClientHandshake(t *testing.T) {
	conn, err := client.Handshake("ws://127.0.0.1:10000/")
	if err != nil {
		t.Fatalf("Error during handshake: %v", err)
	}

	message := []byte("ping!")
	f := client.WSFrame{
		Fin:        false,
		Opcode:     0x9,
		Mask:       true,
		PayloadLen: uint64(len(message)),
		Payload:    message,
	}
	resp, err := sendFrame(conn, f)
	if err != nil {
		t.Fatalf("Error during handshake: %v", err)
	}
	decodedResp, _ := client.Decode(resp)
	log.Printf("Decoded Response: %v", decodedResp)

	f = client.WSFrame{
		Fin:        true,
		Opcode:     0x8,
		Mask:       true,
		PayloadLen: uint64(len(message)),
		Payload:    message,
	}
	_, err = sendFrame(conn, f)
	if err != nil {
		t.Fatalf("Error sending frame: %v", err)
	}
}
