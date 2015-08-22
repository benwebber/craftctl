package config

import (
	"net"

	"github.com/codegangsta/cli"
)

type Config struct {
	Address  *net.TCPAddr
	Password string
}

func New() Config {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:25575")
	return Config{
		Address:  addr,
		Password: "password",
	}
}

func NewConfigFromContext(ctx *cli.Context) (Config, error) {
	host := ctx.String("host")
	port := ctx.String("port")
	password := ctx.String("password")
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
	return Config{
		Address:  addr,
		Password: password,
	}, err
}
