package rcon

import (
	"bytes"
	"strings"
)

type Request struct {
	*Packet
}

func NewRequest(t packetType, payload ...string) (*Request, error) {
	var buf bytes.Buffer
	// Validate request packet type.
	switch t {
	case SERVERDATA_EXECCOMMAND:
		// Prepend with / if necessary.
		if !strings.HasPrefix(payload[0], "/") {
			buf.WriteString("/")
		}
	case SERVERDATA_AUTH:
		fallthrough
	default:
		return &Request{}, InvalidPacketType
	}
	// Separate command fragments with spaces.
	for i, fragment := range payload {
		if i != 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fragment)
	}
	return &Request{
		NewPacket(t, buf.String()),
	}, nil
}

func NewCommand(payload ...string) (*Request, error) {
	return NewRequest(SERVERDATA_EXECCOMMAND, payload...)
}
