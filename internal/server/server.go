package server

import (
	"bytes"
	"errors"
	"go_redis/pkg/goredis"
	"log"
	"net"
	"sync"
)

type CacheServer struct {
	listenAddr     string
	listener       net.Listener
	quitChannel    chan struct{}
	msgChannel     chan Message
	addPeerChannel chan *net.Conn
	peers          map[*Peer]bool
	cache          Cache
	wg             *sync.WaitGroup
}

func NewCacheServer(listenAddr string) *CacheServer {
	return &CacheServer{
		listenAddr:     listenAddr,
		quitChannel:    make(chan struct{}),
		msgChannel:     make(chan Message),
		addPeerChannel: make(chan *net.Conn),
		peers:          make(map[*Peer]bool),
		cache:          NewLRUCache(4),
		wg:             new(sync.WaitGroup),
	}
}

func (s *CacheServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Println("Error starting the listener:", err.Error())
		return
	}
	s.listener = listener

	s.wg.Add(2)
	go s.acceptConnections()
	go s.mainLoop()

	s.wg.Wait()
}

func (s *CacheServer) mainLoop() {
	defer s.wg.Done()
	for {
		select {
		case <-s.quitChannel:
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
		err = msg.peer.WriteData(bytes.NewBufferString(s.cache.Get(key)).Bytes())
		if err != nil {
			return err
		}
		return nil
	case PutCommand:
		key := cmd.(PutCommand).key
		val := cmd.(PutCommand).val
		s.cache.Put(key, val)
		err = msg.peer.WriteData(bytes.NewBufferString("DONE").Bytes())
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid message format")
}

func (s *CacheServer) attachPeer(conn *net.Conn) {
	peer := NewPeer(*conn, s.msgChannel)
	s.peers[peer] = true
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := peer.ReadData()
		if errors.Is(err, goredis.CONN_CLOSE) {
			return
		}
		if err != nil {
			log.Println("Error while reading data:", err.Error())
		}
		s.peers[peer] = false
	}()
}

func (s *CacheServer) acceptConnections() {
	defer s.wg.Done()
	for {
		select {
		case <-s.quitChannel:
			return
		default:
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
}

func (s *CacheServer) Quit() {
	var wg sync.WaitGroup
	for peer, connected := range s.peers {
		if !connected {
			continue
		}
		wg.Add(1)
		peer.SendCloseMessage(&wg)
	}
	wg.Wait()
	for peer, connected := range s.peers {
		if !connected {
			continue
		}
		wg.Add(1)
		peer.Close(&wg)
	}
	wg.Wait()
	s.quitChannel <- struct{}{}
	log.Println("Closing TCP listener...")
	err := s.listener.Close()
	if err != nil {
		log.Println("Error closing the listener:", err.Error())
	}
}
