package goredis

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var CONN_CLOSE = errors.New("CONN_CLOSE")

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
	if sizeInt == 0 {
		return nil, CONN_CLOSE
	}
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
