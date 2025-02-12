package main

import (
	"log"
	"net/http"
)

func do_handshake(client *http.Client) (*http.Client, error) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:9001/chat", nil)
	req.Header.Add("Host", "my.websocket.com")
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "Upgrade")
	// What's up with the key???
	req.Header.Add("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	// What are these modes :0
	req.Header.Add("Sec-WebSocket-Protocol", "chat, superchat")
	req.Header.Add("Sec-WebSocket-Version", "13")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching url:", err)
		return nil, err
	}
	log.Println("Initialized client...")
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
	do_handshake(client)
}
