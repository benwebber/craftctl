package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
)

type packetType int32

// Packet represents an RCON request or response packet.
type Packet struct {
	Size    int32      // Size of message
	Id      int32      // Message ID
	Type    packetType // Message type (see SERVERDATA_ constants)
	Payload []byte     // Message payload
}

// NewPacket constructs a new Packet of type t that delivers payload p.
func NewPacket(t packetType, payload string) *Packet {
	// ID (int32) + type (int32) + null-terminated string + null string
	size := 2*binary.Size(int32(0)) + 2*len(NUL) + len(payload)
	return &Packet{
		Size:    int32(size),
		Id:      rand.Int31(),
		Type:    t,
		Payload: []byte(payload),
	}
}

// String implements the Stringer interface for Packet.
func (p *Packet) String() string {
	return fmt.Sprintf("%x %x %x %x", p.Size, p.Id, p.Type, p.Payload)
}

// Marshall converts the Packet into its binary representation.
func (r Packet) Marshall() ([]byte, error) {
	buf := new(bytes.Buffer)
	rw := &binaryErrorReadWriter{w: buf, byteOrder: binary.LittleEndian}
	rw.write(r.Size)
	rw.write(r.Id)
	rw.write(r.Type)
	// null-terminate command and packet
	payload := append(r.Payload, 0x0, 0x0)
	rw.write(payload)
	if rw.err != nil {
		return []byte{}, rw.err
	}
	return buf.Bytes(), nil
}

// Unmarshall unpacks the binary representation of a Packet.
func Unmarshall(data []byte) (Packet, error) {
	size, err := byteSliceToInt32(data[0:4])
	id, err := byteSliceToInt32(data[4:8])
	type_, err := byteSliceToInt32(data[8:12])
	return Packet{
		Size:    size,
		Id:      id,
		Type:    packetType(type_),
		Payload: bytes.Trim(data[12:], NUL),
	}, err
}
