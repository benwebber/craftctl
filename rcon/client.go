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
func (c *Client) Auth() (*Response, error) {
	r, err := NewAuthRequest(c.Config.Password)
	if err != nil {
		return &Response{}, err
	}
	return c.Execute(r)
}

func (c *Client) Execute(r *Request) (*Response, error) {
	c.write(r.Packet)
	resp, err := c.read()
	if err != nil {
		return resp, err
	}
	if resp.Id == int32(SERVERDATA_AUTH_FAILED) {
		err = AuthenticationError
	}
	return resp, err
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
func (c *Client) read() (*Response, error) {
	buf := make([]byte, 4096)
	_, err := c.Conn.Read(buf)
	if err != nil {
		return &Response{}, err
	}
	p, err := Unmarshall(buf)
	return &Response{Packet: &p}, err
}
