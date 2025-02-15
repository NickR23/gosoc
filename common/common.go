package client

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
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

	/**
			The RSV{1,2,3} bits should all be set to zero in this first byte
			"unless an extension is negotiated that defines meanings
	      for non-zero values"
		**/

	log.Println(f.PayloadLen)
	log.Printf("payload: %v", f.Payload)
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

func Decode(frameData []byte) (*WSFrame, error) {
	if len(frameData) < 2 {
		return nil, errors.New("frame too short")
	}
	frame := &WSFrame{}
	frame.Fin = (frameData[0] & 0x80) != 0
	frame.Opcode = frameData[0] & 0x0f
	payloadLen := uint64(frameData[1] & 0x7F) // opcodes are 7 bits wide
	offset := 2                               //Start reading from the 2nd byte

	if payloadLen == 126 { // Extended payload
		if len(frameData) < 4 {
			return nil, errors.New("invalid frame: missing extended payload lentgth")
		}
		payloadLen = uint64(binary.BigEndian.Uint16(frameData[2:4]))
	} else if payloadLen == 127 { // Huge payload length (8 bytes wide)
		if len(frameData) < 10 {
			return nil, errors.New("invalid frame: missing extended payload length")
		}
		payloadLen = binary.BigEndian.Uint64(frameData[2:10])
	}
	if len(frameData) < int(offset+int(payloadLen)) {
		return nil, errors.New("invalid frame: incomplete payload data")
	}
	frame.Payload = frameData[offset : offset+int(payloadLen)]
	return frame, nil
}

func (f *WSFrame) String() string {
	payloadStr := string(f.Payload)
	if !isPrintable(f.Payload) {
		payloadStr = hex.EncodeToString(f.Payload)
	}
	return fmt.Sprintf(
		"WSFrame(FIN=%t, Opcode=0x%X, PayloadLen=%d, Payload=%s)",
		f.Fin, f.Opcode, f.PayloadLen, payloadStr)
}

func isPrintable(data []byte) bool {
	for _, b := range data {
		if b < 32 || b > 126 {
			return false
		}
	}
	return true
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
