package storage

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
)

type ISocketStore interface {
	Store(c socketio.Conn) (string, error)
	Retrieve(id string) (socketio.Conn, error)
	Remove(id string) error
}

type SocketStoreInMemory struct {
	sockets []socketio.Conn
}

func NewSocketStoreInMemory() *SocketStoreInMemory{
	return &SocketStoreInMemory{
		sockets: []socketio.Conn{},
	}
}

func (s *SocketStoreInMemory) Store(c socketio.Conn) (string, error) {
	if c == nil {
		return "", fmt.Errorf("SocketStoreInMemory.Store: cannot save a null socket connection")
	}

	s.sockets = append(s.sockets, c)
	return c.ID(), nil
}

func (s *SocketStoreInMemory) Retrieve(id string) (socketio.Conn, error) {
	for _, socket := range s.sockets {
		if socket.ID() == id {
			return socket, nil
		}
	}
	return nil, fmt.Errorf("SocketStoreInMemory.Retrieve: socket connection %s not found", id)
}

func (s *SocketStoreInMemory) Remove(id string) error {
	for index, socket := range s.sockets {
		if socket.ID() == id {
			next := index + 1
			s.sockets = append(s.sockets[:index], s.sockets[next:]...)
			return nil
		}
	}
	return fmt.Errorf("SocketStoreInMemory.Remove: socket connection %s not found", id)
}
