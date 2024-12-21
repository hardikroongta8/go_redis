package server

import (
	"go_redis/pkg/goredis"
	"io"
	"net"
)

type Peer struct {
	conn       net.Conn
	msgChannel chan Message
}

func NewPeer(conn net.Conn, msgChannel chan Message) *Peer {
	return &Peer{
		conn:       conn,
		msgChannel: msgChannel,
	}
}

type Message struct {
	data []byte
	peer *Peer
}

func (p *Peer) ReadData() error {
	for {
		parser := goredis.NewReader(p.conn)
		data, err := parser.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		p.msgChannel <- Message{
			data: data,
			peer: p,
		}
	}
}

func (p *Peer) WriteData(data []byte) error {
	w := goredis.NewWriter(p.conn)
	return w.Write(data)
}
