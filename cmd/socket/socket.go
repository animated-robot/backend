package socket

import (
	"animated-robot/storage"
	"animated-robot/tools"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type ISocketFactory interface {
	New() *socketio.Server
}

type SocketFactory struct {
	uuidGenerator tools.UUIDGenerator
	socketStore storage.ISocketStore
	sessionStore storage.ISessionStore
	log *logrus.Logger
}

func NewSocketFactory(socketStore storage.ISocketStore, sessionStore storage.ISessionStore, uuidGenerator tools.UUIDGenerator, log *logrus.Logger) ISocketFactory {
	return SocketFactory{
		uuidGenerator: uuidGenerator,
		socketStore:  socketStore,
		sessionStore: sessionStore,
		log:          log,
	}
}

func (sf SocketFactory) New() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		sf.log.WithFields(logrus.Fields{
			"message": err.Error(),
		}).Fatal("NewSocket: Error creating socket")
	}

	setupNamespace(server, NewFrontNamespace(sf.socketStore, sf.sessionStore, sf.log))
	setupNamespace(server, NewInputNamespace(sf.socketStore, sf.sessionStore, sf.log, sf.uuidGenerator))

	return server
}

func setupNamespace(s *socketio.Server, sn ISocketNamespace) {
	namespace := sn.GetNamespace()
	s.OnConnect(namespace, sn.OnConnect)
	s.OnDisconnect(namespace, sn.OnDisconnect)
	s.OnError(namespace, sn.OnError)
	for k, v := range sn.OnEvents() {
		s.OnEvent(namespace, k, v)
	}
}

type ISocketNamespace interface {
	GetNamespace() string
	OnConnect(s socketio.Conn) error
	OnDisconnect(s socketio.Conn, reason string)
	OnError(s socketio.Conn, e error)
	OnEvents() map[string]OnEventHandler
}

type OnEventHandler func(s socketio.Conn, msg string)
