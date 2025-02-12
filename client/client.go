package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
)

type WSFrame struct {
	Fin        bool
	Opcode     byte
	Mask       bool
	PayloadLen uint64 // This should be compressed when serialized.
	MaskingKey []byte
	Payload    []byte
}

func applyMask(payload []byte, maskKey []byte) []byte {
	maskedPayload := make([]byte, len(payload))
	for i := range payload {
		maskedPayload[i] = payload[i] ^ maskKey[i%4]
	}
	return maskedPayload
}

func generateMaskingKey() []byte {
	maskKey := make([]byte, 4)
	for i := range maskKey {
		maskKey[i] = byte(rand.Intn(256))
	}
	return maskKey
}

func (f *WSFrame) Encode() ([]byte, error) {
	firstByte := byte(0)
	if f.Fin {
		firstByte |= 0x80 // This sets the Fin bit.
	}
	firstByte |= f.Opcode // Lower bits are the opcode i think...

	payloadLenBytes, secondByte := encodePayloadLength(f.PayloadLen)
	if f.Mask {
		secondByte |= 0x80 // Mask bit is the higher bit here
	}

	var maskKey []byte
	payload := f.Payload

	if f.Mask {
		maskKey = generateMaskingKey()
		payload = applyMask(payload, maskKey)
	}

	frame := []byte{firstByte, secondByte}
	frame = append(frame, payloadLenBytes...)
	if f.Mask {
		frame = append(frame, maskKey...)
	}
	frame = append(frame, payload...)
	return frame, nil

}

func encodePayloadLength(payloadLen uint64) ([]byte, byte) {
	var lengthBytes []byte
	secondByte := byte(0)

	switch {
	case payloadLen <= 125:
		secondByte |= byte(payloadLen)
	case payloadLen <= 65535:
		secondByte |= 126
		lengthBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(lengthBytes, uint16(payloadLen))
	default:
		secondByte |= 127
		lengthBytes = make([]byte, 8)
		binary.BigEndian.PutUint64(lengthBytes, payloadLen)
	}
	return lengthBytes, secondByte
}

// Returns an initial handshake for the websocket session
func Handshake(wsURL string) (net.Conn, error) {

	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, err
	}
	log.Printf("path: %v\nhost: %v\n\n", u.Path, u.Host)

	host := u.Host
	clientKey := "dGhlIHNhbXBsZSBub25jZQ=="
	// What's up with the key???
	/**
	   To prove that the handshake was received, the server has to take two
	   pieces of information and combine them to form a response.  The first
	   piece of information comes from the |Sec-WebSocket-Key| header field
	   in the client handshake:

	        Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==

	   For this header field, the server has to take the value (as present
	   in the header field, e.g., the base64-encoded [RFC4648] version minus
	   any leading and trailing whitespace) and concatenate this with the
	   Globally Unique Identifier (GUID, [RFC4122]) "258EAFA5-E914-47DA-
	   95CA-C5AB0DC85B11" in string form, which is unlikely to be used by
	   network endpoints that do not understand the WebSocket Protocol.  A
	   SHA-1 hash (160 bits) [FIPS.180-3], base64-encoded (see Section 4 of
	   [RFC4648]), of this concatenation is then returned in the server's
	   handshake.
		**/
	req := fmt.Sprintf(
		"GET %s HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Key: %s\r\n"+
			"Sec-WebSocket-Version: 13\r\n\r\n",
		u.Path, u.Host, clientKey,
	)
	log.Printf("Request: %v", req)

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write([]byte(req))
	if err != nil {
		return nil, err
	}

	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, err
	}

	responseStr := string(resp[:n])
	log.Println("Server Handshake response:", responseStr)
	return conn, nil
}

// Initialize an http client for websockets.
// Starts the websocket handshake.
func Initialize(client *http.Client) *http.Client {
	log.Println("Initializing client...")
	// Space for client initialization.
	// Put all of your policy crap here
	log.Println("Initialized client...")
	return client
}

func main() {
}
