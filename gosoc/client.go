package main

import (
	"log"
	"net/http"
)

// Returns an initial handshake for the websocket session
// TODO some of these keys are configurable. Let's not hardcode
func build_handshake() (*http.Request, error) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:9001/chat", nil)
	req.Header.Add("Host", "my.websocket.com")
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "Upgrade")
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
	req.Header.Add("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	// What are these modes :0
	req.Header.Add("Sec-WebSocket-Protocol", "chat, superchat")
	req.Header.Add("Sec-WebSocket-Version", "13")
	if err != nil {
		return nil, err
	}
	return req, nil

}

func do_handshake(client *http.Client, req *http.Request) (*http.Client, error) {
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching url:", err)
		return nil, err
	}
	log.Println("Response:", resp.Status)
	return client, err

}

// Initialize an http client for websockets.
// Starts the websocket handshake.
func initialize(client *http.Client) *http.Client {
	log.Println("Initializing client...")
	// Space for client initialization.
	// Put all of your policy crap here
	log.Println("Initialized client...")
	return client
}

func main() {
	log.Println("Running client...")
	client := &http.Client{}
	initialize(client)
	handshake_req, _ := build_handshake()
	resp, err := client.Do(handshake_req)
	if err != nil {
		log.Println("Error fetching url:", err)
	}
	log.Println("Response:", resp.Status)
}
