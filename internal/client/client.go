package client

import (
	"context"
	"fmt"
	"go_redis/pkg/goredis"
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
		reader := goredis.NewReader(c.conn)
		data, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println(string(data))
	}
}

func (c *Client) Put(ctx context.Context, key string, value string) error {
	str := fmt.Sprintf("PUT\r\n%s\r\n%s\r\n", key, value)
	writer := goredis.NewWriter(c.conn)
	return writer.Write([]byte(str))
}

func (c *Client) Get(ctx context.Context, key string) error {
	str := fmt.Sprintf("GET\r\n%s\r\n", key)
	writer := goredis.NewWriter(c.conn)
	return writer.Write([]byte(str))
}

func (c *Client) Close() error {
	return c.conn.Close()
}
