package rcon

import (
	"net"

	"github.com/benwebber/craftctl/config"
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
func (c *Client) Auth() *Response {
	r, err := NewRequest(SERVERDATA_AUTH, c.Config.Password)
	if err != nil {
		return &Response{
			Packet: &Packet{},
			Err:    err,
		}
	}
	return c.Execute(r)
}

func (c *Client) Execute(r *Request) *Response {
	c.write(r.Packet)
	resp := c.read()
	if r.Type == SERVERDATA_AUTH && resp.Id == int32(SERVERDATA_AUTH_FAILED) {
		resp.Err = AuthenticationError
	}
	return resp
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
func (c *Client) read() *Response {
	buf := make([]byte, 4096)
	_, err := c.Conn.Read(buf)
	if err != nil {
		return &Response{
			Packet: &Packet{},
			Err:    err,
		}
	}
	p, err := Unmarshall(buf)
	return &Response{
		Packet: &p,
		Err:    err,
	}
}
