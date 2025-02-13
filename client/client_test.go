package client_test

import (
	"log"
	"net"
	"testing"

	client "github.com/NickR23/gosoc/client"
	common "github.com/NickR23/gosoc/common"
)

func sendFrame(conn net.Conn, f common.WSFrame) ([]byte, error) {
	encodedFrame, _ := f.Encode()
	decodedFrame, _ := common.Decode(encodedFrame)
	log.Printf("Sending frame: %v", decodedFrame)
	_, err := conn.Write(encodedFrame)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Example test that runs against the server
func TestClientHandshake(t *testing.T) {
	conn, err := client.Handshake("ws://127.0.0.1:10000/")
	if err != nil {
		t.Fatalf("Error during handshake: %v", err)
	}

	message := []byte("ping!")
	f := common.WSFrame{
		Fin:        false,
		Opcode:     0x1,
		Mask:       true,
		PayloadLen: uint64(len(message)),
		Payload:    message,
	}
	resp, err := sendFrame(conn, f)
	if err != nil {
		t.Fatalf("Error during handshake: %v", err)
	}
	decodedResp, _ := common.Decode(resp)
	log.Printf("Decoded Response: %v", decodedResp)

	f = common.WSFrame{
		Fin:        true,
		Opcode:     0x8,
		Mask:       true,
		PayloadLen: uint64(len(message)),
		Payload:    message,
	}
	resp, err = sendFrame(conn, f)
	decodedResp, err = common.Decode(resp)
	log.Printf("Decoded response: %v", decodedResp)
	if err != nil {
		t.Fatalf("Error sending frame: %v", err)
	}
}
