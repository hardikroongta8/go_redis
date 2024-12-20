package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr: addr,
		conn: conn,
	}, nil
}

func (c *Client) ReadData() error {
	for {
		buf := make([]byte, 2048)
		n, err := c.conn.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println(string(buf[:n]))
	}
}

func (c *Client) Put(ctx context.Context, key string, value string) error {
	str := fmt.Sprintf("PUT\r\n%s\r\n%s\r\n", key, value)
	_, err := c.conn.Write(bytes.NewBufferString(str).Bytes())
	return err
}

func (c *Client) Get(ctx context.Context, key string) error {
	str := fmt.Sprintf("GET\r\n%s\r\n", key)
	_, err := c.conn.Write(bytes.NewBufferString(str).Bytes())
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
