package server

import (
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
		buf := make([]byte, 1024)
		n, err := p.conn.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		p.msgChannel <- Message{
			data: data,
			peer: p,
		}
	}
}

func (p *Peer) WriteData(data []byte) error {
	_, err := p.conn.Write(data)
	return err
}
