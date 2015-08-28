package rcon

type Response struct {
	*Packet
	Err error
}

func (r *Response) String() string {
	return string(r.Payload)
}
