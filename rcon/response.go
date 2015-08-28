package rcon

type Response struct {
	*Packet
}

func (r *Response) String() string {
	return string(r.Payload)
}
