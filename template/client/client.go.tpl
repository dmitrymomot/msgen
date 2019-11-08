package client

import (
	"fmt"

	"{{ .ServicePath }}/pb/{{ package .ServiceName }}"
	"google.golang.org/grpc"
)

type (
	// Client interface
	Client interface {
		{{ package .ServiceName }}.ServiceClient
		Close() error
	}

	client struct {
		{{ package .ServiceName }}.ServiceClient
		conn *grpc.ClientConn
	}
)

// Close gRPC connection with examplesrv service
func (c *client) Close() error {
	return c.conn.Close()
}

// NewDefaultClient returns grpc connection
func NewDefaultClient(addr string) (cl Client, err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect to %s: %v", addr, err)
	}
	c := {{ package .ServiceName }}.NewServiceClient(conn)
	return &client{c, conn}, nil
}

// Must returns client or panics
func Must(c Client, err error) Client {
	if err != nil {
		panic(err)
	}
	return c
}
