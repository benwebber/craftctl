package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"time"
)

const NUL = "\x00"

// RCON request / response codes
const (
	SERVERDATA_AUTH           packetType = 3
	SERVERDATA_AUTH_RESPONSE  packetType = 2
	SERVERDATA_AUTH_FAILED    packetType = -1
	SERVERDATA_EXECCOMMAND    packetType = 2
	SERVERDATA_RESPONSE_VALUE packetType = 0
)

var (
	AuthenticationError = errors.New("rcon: authentication failed")
	UnknownCommandError = errors.New("rcon: unknown command")
	UsageError          = errors.New("rcon: invalid syntax")
	InvalidJSONError    = errors.New("rcon: invalid JSON")
	PlayerNotFoundError = errors.New("rcon: player not found")
	InvalidPacketType   = errors.New("rcon: invalid packet type")
)

var errorMap = map[string]error{
	"Usage:":                      UsageError,
	"Unknown command.":            UnknownCommandError,
	"Invalid json:":               InvalidJSONError,
	"That player cannot be found": PlayerNotFoundError,
}

// Convert a byte slice into an int32.
func byteSliceToInt32(data []byte) (int32, error) {
	var i int32
	buf := bytes.NewReader(data)
	rw := &binaryErrorReadWriter{r: buf, byteOrder: binary.LittleEndian}
	rw.read(&i)
	return i, rw.err
}

// binaryErrorReadWriter captures errors messages while reading from and
// writing to binary buffers.
type binaryErrorReadWriter struct {
	r         io.Reader
	w         io.Writer
	byteOrder binary.ByteOrder
	err       error
}

func (rw *binaryErrorReadWriter) read(data interface{}) {
	if rw.err != nil {
		return
	}
	rw.err = binary.Read(rw.r, rw.byteOrder, data)
}

func (rw *binaryErrorReadWriter) write(data interface{}) {
	if rw.err != nil {
		return
	}
	rw.err = binary.Write(rw.w, rw.byteOrder, data)
}

// Seed the random number generator for packet IDs.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
