package main

import (
	"log"
	"net/http"
)

func get() (*http.Response, error) {
	resp, err := http.Get("http://google.com")
	if err != nil {
		log.Println("Error fetching url:", err)
		return nil, err
	}
	log.Println("Received response:", resp)
	return resp, nil
}

func initialize(client *http.Client) (*http.Client, error) {
	log.Println("Initializing client...")
	// Is /chat right??????
	req, err := http.NewRequest("GET", "http://127.0.0.1:9001/chat", nil)
	req.Header.Add("Host", "my.websocket.com")
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "Upgrade")
	// What's up with the key???
	req.Header.Add("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
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

func main() {
	log.Println("Running client...")
	resp, _ := get()
	if resp != nil {
		log.Println(resp.Status)
	}

	client := &http.Client{}

	initialize(client)
}
