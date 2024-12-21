package client

import (
	"context"
	"errors"
	"fmt"
	"go_redis/pkg/goredis"
	"log"
	"net"
	"sync"
)

type Client struct {
	addr string
	conn net.Conn
	WG   *sync.WaitGroup
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr: addr,
		conn: conn,
		WG:   new(sync.WaitGroup),
	}, nil
}

func (c *Client) ReadData() error {
	for {
		reader := goredis.NewReader(c.conn)
		data, err := reader.Read()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, goredis.CONN_CLOSE) {
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

func (c *Client) SendCloseMessage(wg *sync.WaitGroup) {
	defer wg.Done()
	w := goredis.NewWriter(c.conn)
	err := w.Write([]byte(""))
	if err != nil {
		log.Println("Error sending CLOSE signal to client:", err.Error())
		return
	}
}

func (c *Client) Close() error {
	writer := goredis.NewWriter(c.conn)
	err := writer.Write([]byte(""))
	if err != nil {
		return err
	}
	return c.conn.Close()
}
