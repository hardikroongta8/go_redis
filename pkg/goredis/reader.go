package goredis

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Reader struct {
	conn net.Conn
}

func NewReader(conn net.Conn) *Reader {
	return &Reader{
		conn: conn,
	}
}

func (p *Reader) Read() ([]byte, error) {
	sizeBytes := make([]byte, 2)
	_, err := p.conn.Read(sizeBytes)
	if err != nil {
		return nil, err
	}
	sizeInt := binary.BigEndian.Uint16(sizeBytes)
	dataBytes := make([]byte, sizeInt)
	n, err := p.conn.Read(dataBytes)
	if err != nil {
		return nil, err
	}
	if n != int(sizeInt) {
		return nil, fmt.Errorf("invalid format of data")
	}
	return dataBytes, nil
}
