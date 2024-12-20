package server

import (
	"bytes"
	"errors"
	"log"
	"net"
)

type CacheServer struct {
	listenAddr     string
	listener       net.Listener
	quitChannel    chan struct{}
	msgChannel     chan Message
	addPeerChannel chan *net.Conn
	peers          []*Peer
	cache          Cache
}

func NewCacheServer(listenAddr string) *CacheServer {
	return &CacheServer{
		listenAddr:     listenAddr,
		quitChannel:    make(chan struct{}),
		msgChannel:     make(chan Message),
		addPeerChannel: make(chan *net.Conn),
		peers:          make([]*Peer, 0),
		cache:          NewLRUCache(4),
	}
}

func (s *CacheServer) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.listener = listener
	go s.acceptConnections()

	s.mainLoop()
	log.Println("Closing TCP listener...")
	return listener.Close()
}

func (s *CacheServer) mainLoop() {
	for {
		select {
		case <-s.quitChannel:
			log.Println("Stopping server...")
			return
		case conn := <-s.addPeerChannel:
			s.attachPeer(conn)
		case msg := <-s.msgChannel:
			if err := s.handleMessage(msg); err != nil {
				log.Println("Error while reading msg:", err.Error())
			}
		}
	}
}

func (s *CacheServer) handleMessage(msg Message) error {
	cmd, err := parseMessage(msg.data)
	if err != nil {
		return err
	}
	switch cmd.(type) {
	case GetCommand:
		key := cmd.(GetCommand).key
		msg.peer.WriteData(bytes.NewBufferString(s.cache.Get(key)).Bytes())
		return nil
	case PutCommand:
		key := cmd.(PutCommand).key
		val := cmd.(PutCommand).val
		s.cache.Put(key, val)
		msg.peer.WriteData(bytes.NewBufferString("DONE").Bytes())
		return nil
	}
	return errors.New("invalid message format")
}

func (s *CacheServer) attachPeer(conn *net.Conn) {
	peer := NewPeer(*conn, s.msgChannel)
	s.peers = append(s.peers, peer)
	go func() {
		err := peer.ReadData()
		if err != nil {
			log.Println("Error while reading data:", err.Error())
		}
	}()
}

func (s *CacheServer) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			log.Println("Error while accepting connection:", err.Error())
			continue
		}
		log.Println("New Connection:", conn.RemoteAddr())
		s.addPeerChannel <- &conn
	}
}
