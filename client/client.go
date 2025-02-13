package client

import (
	"fmt"
	"net"
	"net/url"
)

// Returns an initial handshake for the websocket session
func Handshake(wsURL string) (net.Conn, error) {

	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, err
	}

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

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write([]byte(req))
	if err != nil {
		return nil, err
	}

	resp := make([]byte, 1024)
	_, err = conn.Read(resp)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
