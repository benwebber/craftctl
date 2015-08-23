package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/benwebber/craftctl/config"
)

type packetType int32

const NUL = "\x00"

// RCON request / response codes
const (
	SERVERDATA_AUTH           packetType = 3
	SERVERDATA_AUTH_RESPONSE  packetType = 2
	SERVERDATA_AUTH_FAILED    packetType = -1
	SERVERDATA_EXECCOMMAND    packetType = 2
	SERVERDATA_RESPONSE_VALUE packetType = 0
)

// Client represents an RCON client.
type Client struct {
	Config config.Config
	Conn   *net.TCPConn
}

// NewClient creates a Client and connects to the RCON service.
func NewClient(cfg config.Config) (Client, error) {
	conn, err := net.DialTCP("tcp", nil, cfg.Address)
	if err != nil {
		return Client{}, err
	}
	return Client{
		Config: cfg,
		Conn:   conn,
	}, nil
}

// Authenticate to the RCON service.
func (c *Client) Auth() (err error) {
	_, err = c.send(SERVERDATA_AUTH, c.Config.Password)
	return
}

func (c *Client) Command(command ...string) (resp string, err error) {
	cmd := prepareCommand(command...)
	return c.send(SERVERDATA_EXECCOMMAND, cmd)
}

// Send a packet to the RCON service.
func (c *Client) write(p *Packet) error {
	data, err := p.Marshall()
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(data)
	return err
}

// Read a packet from the RCON service.
func (c *Client) read() (Packet, error) {
	buf := make([]byte, 4096)
	_, err := c.Conn.Read(buf)
	if err != nil {
		return Packet{}, err
	}
	p, err := Unmarshall(buf)
	return p, err
}

// Send a request to the RCON service.
func (c *Client) send(t packetType, command string) (resp string, err error) {
	p := NewPacket(t, command)
	c.write(&p)
	r, err := c.read()
	if err != nil {
		return "", err
	}
	if r.Id == int32(SERVERDATA_AUTH_FAILED) {
		err = fmt.Errorf("authentication failed")
	}
	resp = string(r.Payload)
	return resp, err
}

// Packet represents an RCON request or response packet.
type Packet struct {
	Size    int32      // Size of message
	Id      int32      // Message ID
	Type    packetType // Message type (see SERVERDATA_ constants)
	Payload []byte     // Message payload
}

// NewPacket constructs a new Packet of type t that delivers payload p.
func NewPacket(t packetType, payload string) Packet {
	// ID (int32) + type (int32) + null-terminated string + null string
	size := 2*binary.Size(int32(0)) + 2*len(NUL) + len(payload)
	return Packet{
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

func prepareCommand(command ...string) string {
	var buf bytes.Buffer
	// Prepend with / if necessary.
	if !strings.HasPrefix(command[0], "/") {
		buf.WriteString("/")
	}
	for i, fragment := range command {
		if i != 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fragment)
	}
	return buf.String()
}

// Seed the random number generator for packet IDs.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
